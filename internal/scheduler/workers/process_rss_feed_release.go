package workers

import (
	"log/slog"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/internal/episode/model"
	episodeService "github.com/jusoaresg/gorgon/internal/episode/service"
	filter "github.com/jusoaresg/gorgon/internal/filter"
	filterService "github.com/jusoaresg/gorgon/internal/filter/service"
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

	settings, err := filterService.ResolveSettings(rss.db, show.ID)
	if err != nil {
		return err
	}

	profile, err := filterService.ResolveProfile(rss.db, settings)
	if err != nil {
		return err
	}

	ctx, err := filterService.BuildContext(rss.db, show, ep.Season, ep.Number, settings)
	if err != nil {
		return err
	}

	termKeys := make([]string, 0)
	for _, pattern := range filterService.SearchPatterns(profile) {
		for _, name := range ctx.AllNames() {
			query, err := filter.ExpandQuery(pattern, ctx, name)
			if err != nil {
				continue
			}
			termKeys = append(termKeys, strings.ToLower(query))
		}
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

	processedResponses = episodeService.FilterAndScoreResponses(processedResponses, profile, ctx)
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
	episodeDownloader := episodeService.EpisodeDownloader{}
	episodeDownloader.DownloadEpisode(ep, processedResponses[0])

	return nil
}
