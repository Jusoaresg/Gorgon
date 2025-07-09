package service

import (
	"log/slog"
	"sort"
	"strings"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
)

func FilterRequiredWords(t []schema.SearchResponse) []schema.SearchResponse {
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

func FilterByEpisodeScore(response []schema.SearchResponse) []schema.SearchResponse {
	sort.Slice(response, func(i, j int) bool {
		scoreI := ScoreEpisode(response[i])
		scoreJ := ScoreEpisode(response[j])
		return scoreI > scoreJ
	})
	return response
}
