package service

import (
	"strings"

	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
)

func getPreferredWordsScore(t schema.SearchResponse) int {
	preferredWord := make(map[string]int)
	preferredWord["multisubs"] = 50

	filename := strings.ToLower(t.Filename)

	score := 0
	for word, points := range preferredWord {
		if strings.Contains(filename, word) {
			score += points
		}
	}
	return score
}

func getQualityScore(t schema.SearchResponse) int {
	possibleQuality := make(map[string]int)
	possibleQuality["2560"] = 20
	possibleQuality["1080"] = 20
	possibleQuality["720"] = 10
	possibleQuality["480"] = 0

	for key, quality := range possibleQuality {
		if strings.Contains(t.Filename, key) {
			return quality
		}
	}
	return 0
}

func ScoreEpisode(t schema.SearchResponse) int {
	score := getPreferredWordsScore(t) + getQualityScore(t)

	health := t.Seeders - t.Leechers
	if t.Seeders == 0 {
		health = -50
	} else if health > 30 {
		health = 30
	} else if health < 0 {
		health = 0
	}

	//TODO: Better detection for the filesize
	return score + health
}
