package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpandQuery_ReplacesAliasWithGivenName(t *testing.T) {
	ctx := testContext()

	query, err := ExpandQuery("{alias} S{season:00}E{episode:00}", ctx, "dbz")
	require.NoError(t, err)
	assert.Equal(t, "dbz S01E04", query)

	query, err = ExpandQuery("{alias} S{season:00}E{episode:00}", ctx, "dragon ball")
	require.NoError(t, err)
	assert.Equal(t, "dragon ball S01E04", query)
}

func TestExpandQuery_LiteralTextIsKept(t *testing.T) {
	ctx := testContext()

	query, err := ExpandQuery("EP{episode} 1080p {alias}", ctx, "dbz")
	require.NoError(t, err)
	assert.Equal(t, "EP4 1080p dbz", query)
}

func TestExpandQuery_UnpaddedSeasonEpisode(t *testing.T) {
	ctx := testContext()

	query, err := ExpandQuery("{season}x{episode}", ctx, "dbz")
	require.NoError(t, err)
	assert.Equal(t, "1x4", query)
}

func TestExpandQuery_NoAliasPlaceholderUsesContext(t *testing.T) {
	ctx := testContext()

	query, err := ExpandQuery("S{season:00}E{episode:00}", ctx, "")
	require.NoError(t, err)
	assert.Equal(t, "S01E04", query)
}

func TestExpandQuery_UnknownPlaceholder(t *testing.T) {
	_, err := ExpandQuery("{unknown}", testContext(), "dbz")
	require.Error(t, err)
}
