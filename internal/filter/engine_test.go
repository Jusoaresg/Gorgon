package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testProfile() *Profile {
	return &Profile{
		ID:   1,
		Name: "test",
		Required: []string{
			"{alias} S{season:00}E{episode:00}",
			"multisub",
		},
		Rejected: []string{
			"hdtv",
			"x265",
		},
		Preferred: []PreferredPattern{
			{Pattern: "web", Score: 30},
			{Pattern: "1080p", Score: 20},
			{Pattern: "720p", Score: 10},
		},
	}
}

func TestEvaluate_AllRequiredPresent(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Dragon Ball S01E04 1080p MULTISUB WEB")
	assert.True(t, result.Passed)
	assert.Equal(t, 50, result.PreferredScore)
}

func TestEvaluate_MissingRequired(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Dragon Ball S01E04 1080p WEB")
	assert.False(t, result.Passed)
	assert.Contains(t, result.RejectedReason, "multisub")
}

func TestEvaluate_AliasDoesNotMatchRequired(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Naruto S01E04 MULTISUB 1080p")
	assert.False(t, result.Passed)
	assert.Contains(t, result.RejectedReason, "{alias}")
}

func TestEvaluate_RejectedMatches(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Dragon Ball S01E04 MULTISUB HDTV")
	assert.False(t, result.Passed)
	assert.Contains(t, result.RejectedReason, "hdtv")
}

func TestEvaluate_RejectedWinsOverRequired(t *testing.T) {
	profile := &Profile{
		Required: []string{"x265"},
		Rejected: []string{"x265"},
	}

	result := Evaluate(profile, testContext(), "Dragon Ball S01E04 X265")
	assert.False(t, result.Passed)
	assert.Contains(t, result.RejectedReason, "rejected: x265")
}

func TestEvaluate_PreferredSumsScores(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Dragon Ball S01E04 MULTISUB 1080p WEB")
	assert.True(t, result.Passed)
	assert.Equal(t, 50, result.PreferredScore)
}

func TestEvaluate_PreferredNoMatchScoresZero(t *testing.T) {
	result := Evaluate(testProfile(), testContext(), "Dragon Ball S01E04 MULTISUB")
	assert.True(t, result.Passed)
	assert.Equal(t, 0, result.PreferredScore)
}

func TestEvaluate_NilProfilePasses(t *testing.T) {
	result := Evaluate(nil, testContext(), "anything at all")
	assert.True(t, result.Passed)
	assert.Equal(t, 0, result.PreferredScore)
}

func TestEvaluate_EmptyProfilePasses(t *testing.T) {
	result := Evaluate(&Profile{}, testContext(), "anything at all")
	assert.True(t, result.Passed)
}

func TestEvaluate_RequiredPatternWithAliasContext(t *testing.T) {
	profile := &Profile{
		Required: []string{"{alias} EP{episode:00}"},
	}

	ctx := testContext()
	assert.True(t, Evaluate(profile, ctx, "DBZ EP04 1080p").Passed)
	assert.False(t, Evaluate(profile, ctx, "Naruto EP04 1080p").Passed)
}

func TestValidate_ValidProfile(t *testing.T) {
	err := Validate(testProfile(), testContext())
	assert.NoError(t, err)
}

func TestValidate_InvalidProfile(t *testing.T) {
	profile := &Profile{
		Required: []string{"{unknown} 1080p"},
	}
	err := Validate(profile, testContext())
	assert.Error(t, err)
}
