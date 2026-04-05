package workers

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeService "github.com/jusoaresg/gorgon/internal/episode/service"
	showService "github.com/jusoaresg/gorgon/internal/show/service"
	"github.com/jusoaresg/gorgon/utils"

	seasonRepository "github.com/jusoaresg/gorgon/internal/season/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
)

type RssReleaseProcessor struct {
	db         *sqlx.DB
	showRepo   showRepository.ShowRepositoryInterface
	seasonRepo seasonRepository.SeasonRepositoryInterface
}

func NewRssReleaseProcessor(db *sqlx.DB) RssReleaseProcessor {
	return RssReleaseProcessor{
		db:         db,
		showRepo:   showRepository.NewShowRepository(db),
		seasonRepo: seasonRepository.NewSeasonRepository(db),
	}
}

func (rss *RssReleaseProcessor) RssProcessRelease(ep model.Episode, responses []schema.SearchResponse) error {
	logger := config.GetLogger().WithGroup("worker").With("name", "RssProcessRelease")

	show, err := rss.showRepo.GetById(ep.ShowID)
	if err != nil {
		return err
	}

	aliases, err := showService.GetNormalizedTitleAlias(show)
	if err != nil {
		logger.Error(
			"Error normalizing show aliases",
			slog.String("error", err.Error()),
			slog.Int64("show_id", ep.ShowID),
			slog.Int64("episode_id", ep.ID),
		)
		return err
	}

	termKeys := make([]string, 0, len(aliases))
	for _, alias := range aliases {
		term := fmt.Sprintf("%s S%02dE%02d", alias, ep.Season, ep.Number)
		termKeys = append(termKeys, strings.ToLower(term))
	}

	var processedResponses []schema.SearchResponse

	responseSeen := make(map[string]struct{})
	for _, response := range responses {
		title := utils.NormalizeTitle(response.Title)

		for _, term := range termKeys {
			if strings.Contains(title, term) {
				if _, seen := responseSeen[title]; !seen {
					processedResponses = append(processedResponses, response)
					responseSeen[title] = struct{}{}
				}
				break
			}
		}
	}

	processedResponses = episodeService.FilterRequiredWords(processedResponses)
	processedResponses = episodeService.FilterByEpisodeScore(processedResponses)
	if len(processedResponses) == 0 {
		logger.Debug("No matching releases after filtering")
		return nil
	}

	logger.Info(
		"Downloading episode",
		slog.String("title", processedResponses[0].Title),
		slog.Int64("show_id", ep.ShowID),
		slog.Int64("episode_id", ep.ID),
	)
	episodeService.DownloadEpisode(ep, processedResponses[0])

	return nil
}
