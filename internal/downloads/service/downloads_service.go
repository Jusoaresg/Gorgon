package service

import (
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/jusoaresg/gorgon/external/qbittorrent/schema"
	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
)

const gorgonCategory = "gorgon"

var finishedStates = map[string]bool{
	"uploading":    true,
	"stalledUP":    true,
	"queuedUP":     true,
	"pausedUP":     true,
	"forcedUP":     true,
	"checkingUP":   true,
	"error":        true,
	"missingFiles": true,
	"unknown":      true,
}

type EpisodeInfo struct {
	EpisodeID   int64  `db:"episode_id"`
	ShowID      int64  `db:"show_id"`
	Name        string `db:"name"`
	Season      int    `db:"season"`
	Number      int    `db:"number"`
	Tracking    string `db:"tracking"`
	ShowName    string `db:"show_name"`
	ShowImage   string `db:"show_image"`
	TorrentHash string `db:"torrent_hash"`
}

type DownloadItem struct {
	Torrent schema.CheckTorrentResponse
	Episode *EpisodeInfo
}

type DownloadsService struct {
	db     *sqlx.DB
	logger *slog.Logger
}

func NewDownloadsService(db *sqlx.DB, logger *slog.Logger) *DownloadsService {
	return &DownloadsService{
		db:     db,
		logger: logger,
	}
}

func (s *DownloadsService) BuildDownloads(torrents []schema.CheckTorrentResponse) ([]DownloadItem, error) {
	episodesByHash, err := s.episodesByHash()
	if err != nil {
		return nil, fmt.Errorf("failed to load episodes by torrent hash: %w", err)
	}

	items := make([]DownloadItem, 0, len(torrents))
	for _, torrent := range torrents {
		if torrent.Category != gorgonCategory {
			continue
		}

		var episode *EpisodeInfo
		if info, ok := episodesByHash[strings.ToLower(torrent.Hash)]; ok {
			ep := info
			episode = &ep
		}

		if finishedStates[torrent.State] {
			if episode == nil || episode.Tracking != episodeModel.TrackingSnatched {
				continue
			}
		}

		items = append(items, DownloadItem{
			Torrent: torrent,
			Episode: episode,
		})
	}

	sortDownloads(items)
	return items, nil
}

func sortDownloads(items []DownloadItem) {
	sort.SliceStable(items, func(i, j int) bool {
		iActive := !finishedStates[items[i].Torrent.State]
		jActive := !finishedStates[items[j].Torrent.State]
		if iActive != jActive {
			return iActive
		}
		return items[i].Torrent.AddedOn < items[j].Torrent.AddedOn
	})
}

func (s *DownloadsService) episodesByHash() (map[string]EpisodeInfo, error) {
	var rows []EpisodeInfo
	query := `
		SELECT e.id AS episode_id, e.show_id, e.name, e.season, e.number, e.tracking, e.torrent_hash,
		       s.name AS show_name, s.image_medium AS show_image
		FROM episodes e
		JOIN shows s ON e.show_id = s.id
		WHERE e.torrent_hash IS NOT NULL AND e.torrent_hash != ''
	`
	if err := s.db.Select(&rows, query); err != nil {
		return nil, err
	}

	episodesByHash := make(map[string]EpisodeInfo, len(rows))
	for _, row := range rows {
		episodesByHash[strings.ToLower(row.TorrentHash)] = row
	}

	return episodesByHash, nil
}
