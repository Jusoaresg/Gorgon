package service

import (
	"log/slog"
	"sort"

	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/external/prowlarr/schema"
	"github.com/jusoaresg/gorgon/internal/filter"
)

type scoredResponse struct {
	response schema.SearchResponse
	score    int
}

// FilterAndScoreResponses applies the profile gates to every response,
// dropping rejected releases and sorting the survivors by score (preferred
// words + quality + health) descending.
func FilterAndScoreResponses(responses []schema.SearchResponse, profile *filter.Profile, ctx filter.Context) []schema.SearchResponse {
	logger := config.GetLogger()

	var scored []scoredResponse
	for _, response := range responses {
		result := filter.Evaluate(profile, ctx, response.Filename)
		if !result.Passed {
			logger.Debug(
				"Release rejected by filter",
				slog.String("filename", response.Filename),
				slog.String("reason", result.RejectedReason),
			)
			continue
		}

		scored = append(scored, scoredResponse{
			response: response,
			score:    result.PreferredScore + baseScore(response),
		})
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	filtered := make([]schema.SearchResponse, 0, len(scored))
	for _, s := range scored {
		filtered = append(filtered, s.response)
	}
	return filtered
}

// IsGoodResponse reports whether a response cleared the profile gates and
// scored high enough to stop searching extra aliases.
func IsGoodResponse(response schema.SearchResponse, profile *filter.Profile, ctx filter.Context) bool {
	const goodScore = 60

	result := filter.Evaluate(profile, ctx, response.Filename)
	if !result.Passed {
		return false
	}

	return result.PreferredScore+baseScore(response) >= goodScore
}
