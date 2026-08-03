package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	showSearchPatternsRepository "github.com/jusoaresg/gorgon/internal/show_search_patterns/repository"
	showSettingsModel "github.com/jusoaresg/gorgon/internal/show_settings/model"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newShowSettingsTestHandler() (*Handler, showSettingsRepository.ShowSettingsRepositoryInterface, showSearchPatternsRepository.ShowSearchPatternsRepositoryInterface, *showRepository.ShowRepository) {
	db := testutils.GetTestDB()
	settingsRepo := showSettingsRepository.NewShowSettingsRepository(db)
	searchPatternsRepo := showSearchPatternsRepository.NewShowSearchPatternsRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	h := &Handler{
		ShowSettingsRepo:       &settingsRepo,
		ShowSearchPatternsRepo: &searchPatternsRepo,
		ShowRepo:               showRepo,
		Logger:                 slog.Default(),
	}
	return h, &settingsRepo, &searchPatternsRepo, showRepo
}

func TestGetShowSettings_DefaultsWhenNoRow(t *testing.T) {
	h, _, _, _ := newShowSettingsTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show-settings/9999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	err := h.GetShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "use_aliases")
	assert.Contains(t, rec.Body.String(), "only_latin")
}

func TestGetShowSettings_Found(t *testing.T) {
	h, settingsRepo, _, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	useAliases := false
	require.NoError(t, settingsRepo.Upsert(showSettingsModel.ShowSettings{
		ShowID:     showID,
		UseAliases: useAliases,
		OnlyLatin:  true,
	}))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show-settings/%d", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.GetShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), `"use_aliases": false`)
}

func TestUpdateShowSettings_Success(t *testing.T) {
	h, settingsRepo, searchPatternsRepo, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	requestJSON := `{"filter_profile_id": null, "use_aliases": false, "only_latin": false, "search_patterns": ["{alias} 4k", "{alias} UHD"]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/database/show-settings/%d", showID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.UpdateShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	settings, err := settingsRepo.GetByShowID(showID)
	require.NoError(t, err)
	assert.False(t, settings.UseAliases)
	assert.False(t, settings.OnlyLatin)

	patterns, err := searchPatternsRepo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Equal(t, []string{"{alias} 4k", "{alias} UHD"}, patterns)
}

func TestUpdateShowSettings_GetReturnsSearchPatterns(t *testing.T) {
	h, settingsRepo, searchPatternsRepo, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	require.NoError(t, settingsRepo.Upsert(showSettingsModel.ShowSettings{
		ShowID:     showID,
		UseAliases: true,
		OnlyLatin:  true,
	}))
	require.NoError(t, searchPatternsRepo.Replace(showID, []string{"{alias} 720p"}))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show-settings/%d", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.GetShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), `"search_patterns"`)
	assert.Contains(t, rec.Body.String(), `"{alias} 720p"`)
}

func TestUpdateShowSettings_EmptyPatternsCleared(t *testing.T) {
	h, _, searchPatternsRepo, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	require.NoError(t, searchPatternsRepo.Replace(showID, []string{"{alias} 4k"}))

	requestJSON := `{"filter_profile_id": null, "use_aliases": true, "only_latin": true, "search_patterns": []}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/database/show-settings/%d", showID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.UpdateShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	patterns, err := searchPatternsRepo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Empty(t, patterns)
}

func TestUpdateShowSettings_EmptyStringsIgnored(t *testing.T) {
	h, _, searchPatternsRepo, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	requestJSON := `{"filter_profile_id": null, "use_aliases": true, "only_latin": true, "search_patterns": ["  ", "{alias} 4k"]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/database/show-settings/%d", showID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.UpdateShowSettings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	patterns, err := searchPatternsRepo.GetByShowID(showID)
	require.NoError(t, err)
	assert.Equal(t, []string{"{alias} 4k"}, patterns)
}

func TestUpdateShowSettings_ShowNotFound(t *testing.T) {
	h, _, _, _ := newShowSettingsTestHandler()

	requestJSON := `{"filter_profile_id": null, "use_aliases": true, "only_latin": true}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/database/show-settings/9999", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	err := h.UpdateShowSettings(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}
