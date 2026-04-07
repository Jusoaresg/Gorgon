package service

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	episodeRepository "github.com/jusoaresg/gorgon/internal/episode/repository"
	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
)

type EpisodeSearchInterface interface {
	ProcessSingleEpisode(episodeID int) error
}

type EpisodeSearchService struct {
	db          *sqlx.DB
	logger      *slog.Logger
	Searcher    EpisodeSearcherInterface
	Downloader  EpisodeDownloaderInterface
	EpisodeRepo episodeRepository.EpisodeRepositoryInterface
	ShowRepo    showRepository.ShowRepositoryInterface
}

func NewEpisodeSearchService(
	db *sqlx.DB,
	logger *slog.Logger,
	searcher EpisodeSearcherInterface,
	downloader EpisodeDownloaderInterface,
	episodeRepo episodeRepository.EpisodeRepositoryInterface,
	showRepo showRepository.ShowRepositoryInterface,
) *EpisodeSearchService {
	return &EpisodeSearchService{
		db:          db,
		logger:      logger,
		Searcher:    searcher,
		Downloader:  downloader,
		EpisodeRepo: episodeRepo,
		ShowRepo:    showRepo,
	}
}

func (s *EpisodeSearchService) ProcessSingleEpisode(episodeID int) error {
	episode, err := s.EpisodeRepo.GetByID(int64(episodeID))
	if err != nil {
		return err
	}

	if aired := episode.HasAired(); !aired {
		s.logger.Info(
			"Episode has not aired yet",
			slog.Int64("episode_id", episode.ID),
			slog.Int64("show_id", episode.ShowID),
		)
		return nil
	}

	show, err := s.ShowRepo.GetById(int64(episode.ShowID))
	if err != nil {
		return err
	}

	s.logger.Info(
		"Searching if episode is available",
		slog.String("show_name", show.Name),
		slog.Int64("show_id", show.ID),
		slog.Int("episode", episode.Number),
		slog.String("episode_name", episode.Name),
		slog.String("tracking", string(episode.Tracking)),
	)

	responses, err := s.Searcher.SearchEpisodeAliasesById(episode, show, s.db)
	if err != nil {
		return err
	}

	if len(responses) == 0 {
		s.logger.Info("No avaible episode found", slog.String("show", show.Name), slog.Int("episode", episode.Number))
		return nil
	}

	responses = FilterRequiredWords(responses)
	if len(responses) <= 0 {
		s.logger.Info("No avaible episode found", slog.String("show", show.Name), slog.Int("episode", episode.Number))
		return nil
	}

	responses = FilterByEpisodeScore(responses)
	if len(responses) <= 0 {
		s.logger.Info(
			"No available episodes found after filtering by episode score",
			slog.Int64("episode_id", episode.ID),
			slog.Int64("show_id", episode.ShowID),
		)
		return nil
	}

	s.logger.Debug("Final chosen torrent", slog.String("filename", responses[0].Filename))
	if err := s.Downloader.DownloadEpisode(episode, responses[0]); err != nil {
		return err
	}
	s.logger.Info("Added torrent to qbittorrent", slog.String("show", show.Name), slog.Int("episode", episode.Number))

	return nil
}
