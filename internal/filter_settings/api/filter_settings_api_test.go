package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	filterSettingsModel "github.com/jusoaresg/gorgon/internal/filter_settings/model"
	filterSettingsRepository "github.com/jusoaresg/gorgon/internal/filter_settings/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFilterSettingsTestHandler() (*Handler, filterSettingsRepository.FilterSettingsRepositoryInterface) {
	db := testutils.GetTestDB()
	repo := filterSettingsRepository.NewFilterSettingsRepository(db)
	h := &Handler{
		FilterSettingsRepo: &repo,
		Logger:             slog.Default(),
	}
	return h, &repo
}

func TestGetFilterSettings_Defaults(t *testing.T) {
	h, _ := newFilterSettingsTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/filter-settings", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetFilterSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), `"use_aliases": true`)
	assert.Contains(t, rec.Body.String(), `"only_latin": true`)
	assert.NotContains(t, rec.Body.String(), `"use_aliases": false`)
}

func TestUpdateFilterSettings_Success(t *testing.T) {
	h, repo := newFilterSettingsTestHandler()

	requestJSON := `{"use_aliases": false, "only_latin": false}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/database/filter-settings", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateFilterSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	settings, err := repo.Get()
	require.NoError(t, err)
	assert.False(t, settings.UseAliases)
	assert.False(t, settings.OnlyLatin)
}

func TestUpdateFilterSettings_PartialTogglesKeepDefault(t *testing.T) {
	h, repo := newFilterSettingsTestHandler()

	profileID := int64(3)
	require.NoError(t, repo.Save(filterSettingsModel.FilterSettings{DefaultFilterProfileID: &profileID, UseAliases: true, OnlyLatin: true}))

	requestJSON := `{"use_aliases": false, "only_latin": false}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/database/filter-settings", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateFilterSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	settings, err := repo.Get()
	require.NoError(t, err)
	assert.False(t, settings.UseAliases)
	assert.False(t, settings.OnlyLatin)
	require.NotNil(t, settings.DefaultFilterProfileID)
	assert.Equal(t, profileID, *settings.DefaultFilterProfileID)
}

func TestUpdateFilterSettings_PartialDefaultKeepsToggles(t *testing.T) {
	h, repo := newFilterSettingsTestHandler()

	requestJSON := `{"default_filter_profile_id": "5"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/database/filter-settings", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateFilterSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	settings, err := repo.Get()
	require.NoError(t, err)
	require.NotNil(t, settings.DefaultFilterProfileID)
	assert.Equal(t, int64(5), *settings.DefaultFilterProfileID)
	assert.True(t, settings.UseAliases)
	assert.True(t, settings.OnlyLatin)
}

func TestUpdateFilterSettings_ClearDefault(t *testing.T) {
	h, repo := newFilterSettingsTestHandler()

	profileID := int64(3)
	require.NoError(t, repo.Save(filterSettingsModel.FilterSettings{DefaultFilterProfileID: &profileID, UseAliases: true, OnlyLatin: true}))

	requestJSON := `{"default_filter_profile_id": ""}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/database/filter-settings", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateFilterSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	settings, err := repo.Get()
	require.NoError(t, err)
	assert.Nil(t, settings.DefaultFilterProfileID)
	assert.True(t, settings.UseAliases)
	assert.True(t, settings.OnlyLatin)
}

func TestUpdateFilterSettings_InvalidDefault(t *testing.T) {
	h, _ := newFilterSettingsTestHandler()

	requestJSON := `{"default_filter_profile_id": "abc"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/database/filter-settings", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateFilterSettings(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}
