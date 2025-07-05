package jobs

import (
	"fmt"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	prowlarr "github.com/jusoaresg/gorgon/external/prowlarr/service"
	qbittorrent "github.com/jusoaresg/gorgon/external/qbittorrent/service"
	"github.com/jusoaresg/gorgon/internal/episode/events"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/utils"
	"log/slog"
	"sort"
	"strings"
	"sync"
	"time"
)

var cooldownCache sync.Map

const cooldownDuration = 5 * time.Minute

func init() {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			now := time.Now()
			cooldownCache.Range(func(key, value any) bool {
				if t, ok := value.(time.Time); ok && now.Sub(t) > cooldownDuration*2 {
					cooldownCache.Delete(key)
				}
				return true
			})
		}
	}()
}

func ProcessSingleEpisode(ep *model.Episode, prowlarrService *prowlarr.ProwlarrSearchService, qbittorrentService *qbittorrent.QBittorrentService) error {
	logger := config.GetLogger().WithGroup("jobs").With("name", "ProcessSingleEpisode")
	db := config.GetSQLite()

	episodeRepo := episodeRepository.NewEpisodeRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	showAliasRepo := showAliasRepository.NewShowAliasesRepository(db)

	show, err := showRepo.GetById(ep.ShowID)
	if err != nil {
		return err
	}

	aired, err := ep.HasAired()
	if err != nil {
		logger.Warn("failed to parse airstamp", slog.Int64("airstamp", ep.AirStamp), slog.String("episode", ep.Name))
	}

	if aired == false {
		return nil
	}

	logger.Info("searching if episode is avaible", slog.String("show_name", show.Name), slog.Int64("show_id", show.ID), slog.Int("episode", ep.Number), slog.String("episode_name", ep.Name), slog.String("tracking", string(ep.Tracking)))

	aliases, err := showAliasRepo.ListByShowID(show.ID)
	titleAlias := []string{
		utils.NormalizeTitle(show.Name),
	}
	for _, alias := range aliases {
		if alias.Alias == show.Name {
			continue
		}
		titleAlias = append(titleAlias, utils.NormalizeTitle(alias.Alias))
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var response []schema.SearchResponse

	for _, title := range titleAlias {
		termKey := fmt.Sprintf("%s S%02dE%02d", title, ep.Season, ep.Number)

		mu.Lock()
		if val, ok := cooldownCache.Load(termKey); ok {
			lastTime := val.(time.Time)
			if time.Since(lastTime) < cooldownDuration {
				logger.Info("skipping search (cooldown)", slog.String("term", termKey))
				mu.Unlock()
				continue
			}
		}

		cooldownCache.Store(termKey, time.Now())
		mu.Unlock()

		wg.Add(1)
		go func(t, key string) {
			defer wg.Done()
			request := schema.SearchRequest{
				Query: key,
			}
			var tmpResponse []schema.SearchResponse
			if err := prowlarrService.Search(&request, &tmpResponse); err != nil {
				logger.Error("error while searching for episodes on prowlarr service", slog.Int("episode", ep.Number), slog.String("show", show.Name))
				return
			}
			mu.Lock()
			response = append(response, tmpResponse...)
			mu.Unlock()
		}(title, termKey)

	}

	wg.Wait()

	if len(response) == 0 {
		logger.Info("no avaible episode found", slog.String("show", show.Name), slog.Int("episode", ep.Number))
		return nil
	}

	response = requiredWords(response)

	if len(response) == 0 {
		logger.Info("no avaible episode found", slog.String("show", show.Name), slog.Int("episode", ep.Number))
	}

	sort.Slice(response, func(i, j int) bool {
		scoreI := score(response[i])
		scoreJ := score(response[j])
		return scoreI > scoreJ
	})

	if len(response) <= 0 {
		//TODO: Log message that was not found any available episode
		return nil
	}

	url := response[0].Guid

	logger.Debug("response", slog.String("show", show.Name), slog.String("episode", ep.Name), slog.String("response", url))
	logger.Debug("final chosen torrent", slog.String("filename", response[0].Filename))

	if err := qbittorrentService.AddTorrent(url); err != nil {
		logger.Error("failed to add torrent", slog.String("error", err.Error()))
		return err
	}

	ep.Tracking = model.TrackingSnatched
	ep.TorrentHash = response[0].InfoHash

	if err := episodeRepo.Update(*ep); err != nil {
		logger.Error("failed to update episode", slog.String("error", err.Error()))
		return err
	}

	logger.Info("added torrent to qbittorrent", slog.String("show", show.Name), slog.Int("episode", ep.Number))

	episode.EmitEpisodeTrackingUpdatedEvent(ep.ID, model.TrackingSnatched)

	return nil
}

func score(t schema.SearchResponse) int {
	//TODO: Better detection for the filesize
	return (t.Seeders - t.Leechers) + preferredWords(t) + getQuality(t) //- int(t.Size)
}

func getQuality(t schema.SearchResponse) int {
	possibleQuality := make(map[string]int)
	possibleQuality["2560"] = 9
	possibleQuality["1080"] = 10
	possibleQuality["720"] = 8
	possibleQuality["480"] = 1

	for key, quality := range possibleQuality {
		if strings.Contains(t.Filename, key) {
			return quality
		}
	}
	return 0
}

func preferredWords(t schema.SearchResponse) int {
	preferredWord := make(map[string]int)
	preferredWord["multisubs"] = 50

	for word, score := range preferredWord {
		if strings.Contains(strings.ToLower(t.Filename), word) {
			return score
		}
	}
	return 0
}

func requiredWords(t []schema.SearchResponse) []schema.SearchResponse {
	logger := config.GetLogger()

	var requiredWords []string
	requiredWords = append(requiredWords, "multisubs")
	requiredWords = append(requiredWords, "multisub")
	requiredWords = append(requiredWords, "multi-sub")
	requiredWords = append(requiredWords, "multi sub")

	var newSchema []schema.SearchResponse
	for _, torrent := range t {
		logger.Debug("checking filename before required words", slog.String("filename", torrent.Filename))
		filename := strings.ToLower(torrent.Filename)
		for _, word := range requiredWords {
			if strings.Contains(filename, word) {
				newSchema = append(newSchema, torrent)
				break
			}
		}
	}
	return newSchema
}
