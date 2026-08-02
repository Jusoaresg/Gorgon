package api

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

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

	requestJSON := `{"default_filter_profile_id": null, "use_aliases": false, "only_latin": false}`
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
