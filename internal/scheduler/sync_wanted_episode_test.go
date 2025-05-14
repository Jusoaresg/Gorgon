package scheduler

import (
	"testing"

	"gorgon/external/prowlarr/schema"
)

func TestRequiredWords(t *testing.T) {
	torrents := []schema.SearchResponse{
		{Filename: "Show.S01E01.MULTISUBS.1080p.WEB-DL.mkv"},
		{Filename: "Show.S01E01.1080p.WEB-DL.mkv"},
		{Filename: "Show.S01E01.multi sub.720p.WEB-DL.mkv"},
		{Filename: "Show.S01E01.multi-sub.720p.WEB-DL.mkv"},
	}

	result := requiredWords(torrents)

	if len(result) != 3 {
		t.Errorf("Expected 3 torrents, got %d", len(result))
	}

	for _, r := range result {
		t.Logf("Accepted: %s", r.Filename)
	}
}
