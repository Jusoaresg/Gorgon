package repository

import (
	"testing"

	"github.com/jusoaresg/gorgon/internal/filter_settings/model"
	"github.com/jusoaresg/gorgon/testutils"
	"github.com/stretchr/testify/assert"
)

func getFilterSettingsRepo() FilterSettingsRepository {
	return NewFilterSettingsRepository(testutils.GetTestDB())
}

func TestFilterSettingsRepository_GetReturnsDefaultsWhenEmpty(t *testing.T) {
	repo := getFilterSettingsRepo()

	got, err := repo.Get()
	assert.NoError(t, err)
	assert.True(t, got.UseAliases)
	assert.True(t, got.OnlyLatin)
	assert.Nil(t, got.DefaultFilterProfileID)
}

func TestFilterSettingsRepository_SaveAndGet(t *testing.T) {
	repo := getFilterSettingsRepo()

	profileID := int64(3)
	settings := model.FilterSettings{
		DefaultFilterProfileID: &profileID,
		UseAliases:             false,
		OnlyLatin:              false,
	}

	assert.NoError(t, repo.Save(settings))

	got, err := repo.Get()
	assert.NoError(t, err)
	assert.NotNil(t, got.DefaultFilterProfileID)
	assert.Equal(t, profileID, *got.DefaultFilterProfileID)
	assert.False(t, got.UseAliases)
	assert.False(t, got.OnlyLatin)
}

func TestFilterSettingsRepository_SaveOverwrites(t *testing.T) {
	repo := getFilterSettingsRepo()

	assert.NoError(t, repo.Save(model.FilterSettings{UseAliases: false, OnlyLatin: true}))
	assert.NoError(t, repo.Save(model.FilterSettings{UseAliases: true, OnlyLatin: true}))

	got, err := repo.Get()
	assert.NoError(t, err)
	assert.True(t, got.UseAliases)
}
