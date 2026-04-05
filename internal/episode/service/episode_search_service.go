package service

import "github.com/jusoaresg/gorgon/internal/episode/repository"

type EpisodeSearchService struct {
	Searcher   EpisodeSearcherInterface
	Downloader EpisodeDownloaderInterface
	Repository repository.EpisodeRepositoryInterface
}

func NewEpisodeSearchService(searcher EpisodeSearcherInterface, downloader EpisodeDownloaderInterface, repo repository.EpisodeRepositoryInterface) *EpisodeSearchService {
	return &EpisodeSearchService{
		Searcher:   searcher,
		Downloader: downloader,
		Repository: repo,
	}
}
