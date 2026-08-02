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
	showSettingsModel "github.com/jusoaresg/gorgon/internal/show_settings/model"
	showSettingsRepository "github.com/jusoaresg/gorgon/internal/show_settings/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newShowSettingsTestHandler() (*Handler, showSettingsRepository.ShowSettingsRepositoryInterface, *showRepository.ShowRepository) {
	db := testutils.GetTestDB()
	settingsRepo := showSettingsRepository.NewShowSettingsRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	h := &Handler{
		ShowSettingsRepo: &settingsRepo,
		ShowRepo:         showRepo,
		Logger:           slog.Default(),
	}
	return h, &settingsRepo, showRepo
}

func TestGetShowSettings_DefaultsWhenNoRow(t *testing.T) {
	h, _, _ := newShowSettingsTestHandler()

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
	h, settingsRepo, showRepo := newShowSettingsTestHandler()

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
	h, settingsRepo, showRepo := newShowSettingsTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	require.NoError(t, err)

	requestJSON := `{"filter_profile_id": null, "use_aliases": false, "only_latin": false}`
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
}

func TestUpdateShowSettings_ShowNotFound(t *testing.T) {
	h, _, _ := newShowSettingsTestHandler()

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
