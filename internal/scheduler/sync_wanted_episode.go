package scheduler

import (
	"fmt"
	"gorgon/config"
	"gorgon/external/prowlarr/schema"
	prowlarr "gorgon/external/prowlarr/service"
	qbittorrent "gorgon/external/qbittorrent/service"
	"gorgon/internal/db/model"
	"gorgon/pkg/services"
	"log/slog"
	"sort"
	"strings"
)

func SyncWantedEpisodes() {

	logger := config.GetLogger()

	baseService := services.NewBaseService()

	var shows []model.Show
	baseService.ListWithPreload(&shows, "Episodes", "Seasons")

	prowlarrService := prowlarr.NewProwlarrSearchService(logger)
	qbittorrentService, err := qbittorrent.NewQBittorrentService(logger)

	if err != nil {
		logger.Error("Error getting qbittorrent service", slog.String("Error", err.Error()))
		return
	}

	for _, show := range shows {
		logger.Info("Searching for new episodes", slog.String("Show", show.Name))

		for _, episode := range show.Episodes {

			aired, err := episode.HasAired()
			if err != nil {
				logger.Warn("Failed to parse AirStamp", slog.String("AirStamp", episode.AirStamp), slog.String("Episode", episode.Name))
				continue
			}

			if aired == false {
				continue
			}

			logger.Debug("Episode status",
				slog.String("Show", show.Name),
				slog.String("Episode", episode.Name),
				slog.String("Tracking", string(episode.Tracking)),
			)
			if episode.Tracking == "wanted" || episode.Tracking == "missing" {
				logger.Info("Searching if episode is avaible", slog.String("Show", show.Name), slog.Int("Episode", episode.Number))

				//TODO: Include titleAlias inside the db

				titleAlias := []string{
					show.Name,
					strings.ReplaceAll(show.Name, "-", ""),
					strings.ReplaceAll(show.Name, "'", ""),
					strings.ReplaceAll(strings.ReplaceAll(show.Name, "'", ""), "-", ""),
				}

				var response []schema.SearchResponse

				for _, title := range titleAlias {
					request := schema.SearchRequest{
						Query: fmt.Sprintf("%s S%02dE%02d", title, episode.Season, episode.Number),
					}

					var tmpResponse []schema.SearchResponse
					if err := prowlarrService.Search(&request, &tmpResponse); err != nil {
						logger.Error("Error while searching for episodes on prowlarrService", slog.Int("Episode", episode.Number), slog.String("Show", show.Name))
						continue
					}
					response = append(response, tmpResponse...)
				}

				if len(response) == 0 {
					logger.Info("No avaible episode found", slog.String("Show", show.Name), slog.Int("Episode", episode.Number))
					continue
				}

				response = requiredWords(response)

				if len(response) == 0 {
					logger.Info("No avaible episode found", slog.String("Show", show.Name), slog.Int("Episode", episode.Number))
					continue
				}

				sort.Slice(response, func(i, j int) bool {
					scoreI := score(response[i])
					scoreJ := score(response[j])
					return scoreI > scoreJ
				})
				url := response[0].Guid

				logger.Debug("Response", slog.String("Show", show.Name), slog.String("Episode", episode.Name), slog.String("Response", url))
				logger.Debug("Final chosen torrent", slog.String("Filename", response[0].Filename))

				qbittorrentService.AddTorrent(url)

				episode.Tracking = model.Tracking.Snatched()
				episode.TorrentHash = response[0].InfoHash
				baseService.UpdateByID(int(episode.ID), &episode)
				logger.Info("Added torrent to qBittorrent", slog.String("Show", show.Name), slog.Int("Episode", episode.Number))
			}
		}
	}
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
