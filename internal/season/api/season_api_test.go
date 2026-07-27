package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	showRepository "github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/internal/season/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newSeasonTestHandler() (*Handler, *repository.SeasonRepository, *showRepository.ShowRepository) {
	db := testutils.GetTestDB()
	seasonRepo := repository.NewSeasonRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	h := &Handler{
		SeasonRepo: seasonRepo,
		Logger:     slog.Default(),
	}
	return h, seasonRepo, showRepo
}

func TestGetShowSeasons_Success(t *testing.T) {
	h, seasonRepo, showRepo := newSeasonTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	season := testutils.MakeFakeSeason()
	season.ShowID = showID
	_, err = seasonRepo.Create(season)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show/season/%d", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.GetShowSeasons(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "Get Show Seasons")
}

func TestGetShowSeasons_InvalidId(t *testing.T) {
	h, _, _ := newSeasonTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/season/abc", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := h.GetShowSeasons(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestGetShowSeasons_Empty(t *testing.T) {
	h, _, _ := newSeasonTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/season/999999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999999")

	err := h.GetShowSeasons(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "Get Show Seasons")
}

func TestGetShowSeasons_MultipleSeasons(t *testing.T) {
	h, seasonRepo, showRepo := newSeasonTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	for i := 0; i < 3; i++ {
		season := testutils.MakeFakeSeason()
		season.ShowID = showID
		season.Number = i + 1
		_, err = seasonRepo.Create(season)
		assert.NoError(t, err)
	}

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show/season/%d", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.GetShowSeasons(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
