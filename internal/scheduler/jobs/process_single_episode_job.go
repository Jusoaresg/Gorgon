package jobs

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	prowlarr "gorgon/external/prowlarr/service"
	qbittorrent "gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/internal/db/repository"
	"gorgon/pkg/handler"
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

func ProcessSingleEpisode(ep *model.Episode, prowlarrService prowlarr.ProwlarrSearchService, qbittorrentService qbittorrent.QBittorrentService) error {
	var logger = config.GetLogger()
	episodeRepo := repository.NewEpisodeRepository()
	showRepo := repository.NewShowRepository()

	show, err := showRepo.GetById(ep.ShowID)
	if err != nil {
		return err
	}

	aired, err := ep.HasAired()
	if err != nil {
		logger.Warn("Failed to parse AirStamp", slog.String("AirStamp", ep.AirStamp), slog.String("Episode", ep.Name))
	}

	if aired == false {
		return nil
	}

	logger.Debug("Episode status",
		slog.String("Show", show.Name),
		slog.String("Episode", ep.Name),
		slog.String("Tracking", string(ep.Tracking)),
	)
	logger.Info("Searching if episode is avaible", slog.String("Show", show.Name), slog.Int("Episode", ep.Number))

	//TODO: Include titleAlias inside the db
	titleAlias := []string{
		normalizeTitle(show.Name),
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
				logger.Info("Skipping search (cooldown)", slog.String("term", termKey))
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
				logger.Error("Error while searching for episodes on prowlarrService", slog.Int("Episode", ep.Number), slog.String("Show", show.Name))
				return
			}
			mu.Lock()
			response = append(response, tmpResponse...)
			mu.Unlock()
		}(title, termKey)

	}

	wg.Wait()

	if len(response) == 0 {
		logger.Info("No avaible episode found", slog.String("Show", show.Name), slog.Int("Episode", ep.Number))
		return nil
	}

	response = requiredWords(response)

	if len(response) == 0 {
		logger.Info("No avaible episode found", slog.String("Show", show.Name), slog.Int("Episode", ep.Number))
	}

	sort.Slice(response, func(i, j int) bool {
		scoreI := score(response[i])
		scoreJ := score(response[j])
		return scoreI > scoreJ
	})
	url := response[0].Guid

	logger.Debug("Response", slog.String("Show", show.Name), slog.String("Episode", ep.Name), slog.String("Response", url))
	logger.Debug("Final chosen torrent", slog.String("Filename", response[0].Filename))

	qbittorrentService.AddTorrent(url)

	ep.Tracking = model.TrackingSnatched
	ep.TorrentHash = response[0].InfoHash

	if err := episodeRepo.Update(*ep); err != nil {
		return err
	}

	logger.Info("Added torrent to qBittorrent", slog.String("Show", show.Name), slog.Int("Episode", ep.Number))

	//TODO: You know
	var msg EpisodeUpdatedWebsocketSchema
	msg.Type = "EpisodeTrackingUpdate"
	msg.EpisodeId = ep.ID
	msg.Tracking = model.TrackingSnatched
	handler.SendWebSocketMessage(msg)
	return nil
}

func normalizeTitle(title string) string {
	replacer := strings.NewReplacer(
		"-", " ", "_", " ", ".", " ",
		"'", "", "\"", "", ":", "", "!", "", ",", "",
	)
	title = replacer.Replace(title)
	title = strings.ToLower(title)
	title = strings.Join(strings.Fields(title), " ") // remove duplicated spaces
	return title
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
		logger.Debug("Checking filename before requiredWords", slog.String("Filename", torrent.Filename))
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
