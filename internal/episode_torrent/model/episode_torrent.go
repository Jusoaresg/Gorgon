package model

import (
	"time"

	prowlarrSchema "github.com/jusoaresg/gorgon/external/prowlarr/schema"
	qbittorrentSchema "github.com/jusoaresg/gorgon/external/qbittorrent/schema"
)

type EpisodeTorrent struct {
	ID        int64 `db:"id"`
	EpisodeId int64 `db:"episode_id"`

	Hash        string `db:"hash"`
	Title       string `db:"title"`
	Indexer     string `db:"indexer"`
	InfoUrl     string `db:"info_url"`
	PublishDate string `db:"publish_date"`
	CreatedAt   int64  `db:"created_at"`
}

func FromSearchResponse(episodeID int64, r prowlarrSchema.SearchResponse) EpisodeTorrent {
	return EpisodeTorrent{
		EpisodeId:   episodeID,
		Hash:        r.InfoHash,
		Title:       r.Title,
		Indexer:     r.Indexer,
		InfoUrl:     r.InfoUrl,
		PublishDate: r.PublishDate,
		CreatedAt:   time.Now().Unix(),
	}
}

func FromAddTorrentRequest(episodeID int64, r qbittorrentSchema.AddEpisodeTorrentRequest) EpisodeTorrent {
	return EpisodeTorrent{
		EpisodeId:   episodeID,
		Hash:        r.InfoHash,
		Title:       r.Title,
		Indexer:     r.Indexer,
		InfoUrl:     r.InfoUrl,
		PublishDate: r.PublishDate,
		CreatedAt:   time.Now().Unix(),
	}
}
