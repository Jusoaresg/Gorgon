package repository

import (
	"testing"

	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getShowSearchPatternsRepo() ShowSearchPatternsRepository {
	return NewShowSearchPatternsRepository(testutils.GetTestDB())
}

func createShowForSearchPatterns(t *testing.T, repo ShowSearchPatternsRepository) int64 {
	t.Helper()
	showRepo := showRepository.NewShowRepository(repo.db)
	id, err := showRepo.Create(testutils.MakeFakeShow())
	require.NoError(t, err)
	return id
}

func TestShowSearchPatternsRepository_ReplaceAndGet(t *testing.T) {
	repo := getShowSearchPatternsRepo()
	showID := createShowForSearchPatterns(t, repo)

	require.NoError(t, repo.Replace(showID, []string{
		"{alias} 4k",
		"{alias} S{season:00}E{episode:00} UHD",
	}))

	patterns, err := repo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Equal(t, []string{"{alias} 4k", "{alias} S{season:00}E{episode:00} UHD"}, patterns)
}

func TestShowSearchPatternsRepository_ReplaceOverwrites(t *testing.T) {
	repo := getShowSearchPatternsRepo()
	showID := createShowForSearchPatterns(t, repo)

	require.NoError(t, repo.Replace(showID, []string{"{alias} 720p"}))
	require.NoError(t, repo.Replace(showID, []string{"{alias} 4k", "{alias} 1080p"}))

	patterns, err := repo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Equal(t, []string{"{alias} 4k", "{alias} 1080p"}, patterns)
}

func TestShowSearchPatternsRepository_ReplaceEmptyClears(t *testing.T) {
	repo := getShowSearchPatternsRepo()
	showID := createShowForSearchPatterns(t, repo)

	require.NoError(t, repo.Replace(showID, []string{"{alias} 4k"}))
	require.NoError(t, repo.Replace(showID, nil))

	patterns, err := repo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Empty(t, patterns)
}

func TestShowSearchPatternsRepository_GetByShowIDEmpty(t *testing.T) {
	repo := getShowSearchPatternsRepo()
	showID := createShowForSearchPatterns(t, repo)

	patterns, err := repo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Empty(t, patterns)
}

func TestShowSearchPatternsRepository_ShowDeletionCascades(t *testing.T) {
	repo := getShowSearchPatternsRepo()
	showID := createShowForSearchPatterns(t, repo)

	require.NoError(t, repo.Replace(showID, []string{"{alias} 4k"}))

	showRepo := showRepository.NewShowRepository(repo.db)
	require.NoError(t, showRepo.DeleteById(showID))

	patterns, err := repo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Empty(t, patterns)
}
