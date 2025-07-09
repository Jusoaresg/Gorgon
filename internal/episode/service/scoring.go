package service

import (
	"strings"

	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
)

func getPreferredWordsScore(t schema.SearchResponse) int {
	preferredWord := make(map[string]int)
	preferredWord["multisubs"] = 50

	for word, score := range preferredWord {
		if strings.Contains(strings.ToLower(t.Filename), word) {
			return score
		}
	}
	return 0
}

func getQualityScore(t schema.SearchResponse) int {
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

func ScoreEpisode(t schema.SearchResponse) int {
	//TODO: Better detection for the filesize
	return (t.Seeders - t.Leechers) + getPreferredWordsScore(t) + getQualityScore(t) //- int(t.Size)
}
