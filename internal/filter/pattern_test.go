package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testContext() Context {
	return Context{
		Show:    "dragon ball",
		Aliases: []string{"dbz", "dragonball"},
		Season:  1,
		Episode: 4,
	}
}

func TestCompile_AliasExpandsToAnyName(t *testing.T) {
	ctx := testContext()
	re, err := Compile("{alias} 1080p", ctx)
	require.NoError(t, err)

	assert.True(t, re.MatchString("Dragon Ball 1080p"))
	assert.True(t, re.MatchString("Dragonball 1080p"))
	assert.True(t, re.MatchString("DBZ 1080p"))
	assert.False(t, re.MatchString("Naruto 1080p"))
}

func TestCompile_ShowOnlyMatchesCanonicalName(t *testing.T) {
	ctx := testContext()
	re, err := Compile("{show} S01E04", ctx)
	require.NoError(t, err)

	assert.True(t, re.MatchString("Dragon Ball S01E04"))
	assert.False(t, re.MatchString("DBZ S01E04"))
}

func TestCompile_SeasonEpisodePadding(t *testing.T) {
	ctx := testContext()

	re, err := Compile("S{season:00}E{episode:00}", ctx)
	require.NoError(t, err)
	assert.True(t, re.MatchString("S01E04"))

	re, err = Compile("S{season}E{episode}", ctx)
	require.NoError(t, err)
	assert.True(t, re.MatchString("S1E4"))
	assert.False(t, re.MatchString("S01E04"))
}

func TestCompile_SeasonEpisodeWiderPadding(t *testing.T) {
	ctx := testContext()

	re, err := Compile("{season:000}{episode:000}", ctx)
	require.NoError(t, err)
	assert.True(t, re.MatchString("001004"))
}

func TestCompile_CaseInsensitive(t *testing.T) {
	ctx := testContext()
	re, err := Compile("{alias} S01e04 1080p", ctx)
	require.NoError(t, err)

	assert.True(t, re.MatchString("dragon ball s01e04 1080p"))
	assert.True(t, re.MatchString("DRAGON BALL S01E04 1080P"))
}

func TestCompile_LiteralRegexCharsAreEscaped(t *testing.T) {
	ctx := testContext()
	re, err := Compile("{alias} 1080p", ctx)
	require.NoError(t, err)

	assert.False(t, re.MatchString("dragon ball x1080xp"))
	assert.True(t, re.MatchString("dragon ball 1080p"))
}

func TestCompile_NoPlaceholdersIsLiteral(t *testing.T) {
	re, err := Compile("multisub", testContext())
	require.NoError(t, err)

	assert.True(t, re.MatchString("My.Show S01E04 MULTISUB 1080p"))
	assert.False(t, re.MatchString("My.Show S01E04 1080p"))
}

func TestCompile_UnknownPlaceholder(t *testing.T) {
	_, err := Compile("{unknown} 1080p", testContext())
	require.Error(t, err)
}

func TestCompile_ReservedPlaceholder(t *testing.T) {
	_, err := Compile("{absolute}", testContext())
	require.Error(t, err)
}

func TestCompile_EmptyPattern(t *testing.T) {
	re, err := Compile("", testContext())
	require.NoError(t, err)
	assert.True(t, re.MatchString("anything"))
}

func TestCompile_AliasOnlyCanonicalWhenNoAliases(t *testing.T) {
	ctx := Context{Show: "one piece", Season: 1, Episode: 2}

	re, err := Compile("{alias} 720p", ctx)
	require.NoError(t, err)
	assert.True(t, re.MatchString("one piece 720p"))
}

func TestCompile_AliasOmitsDuplicates(t *testing.T) {
	ctx := Context{
		Show:    "naruto",
		Aliases: []string{"naruto", "naruto shipudden"},
		Season:  1,
		Episode: 2,
	}

	re, err := Compile("{alias} 1080p", ctx)
	require.NoError(t, err)
	assert.True(t, re.MatchString("naruto 1080p"))
	assert.True(t, re.MatchString("naruto shipudden 1080p"))
}
