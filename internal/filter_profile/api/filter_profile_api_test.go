package api

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	filterProfileModel "github.com/jusoaresg/gorgon/internal/filter_profile/model"
	filterProfileRepository "github.com/jusoaresg/gorgon/internal/filter_profile/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFilterProfileTestHandler() (*Handler, filterProfileRepository.FilterProfileRepositoryInterface) {
	db := testutils.GetTestDB()
	repo := filterProfileRepository.NewFilterProfileRepository(db)
	h := &Handler{
		FilterProfileRepo: &repo,
		Logger:            slog.Default(),
	}
	return h, &repo
}

func TestCreateFilterProfile_Success(t *testing.T) {
	h, repo := newFilterProfileTestHandler()

	requestJSON := `{
		"name": "HD",
		"patterns": [
			{"kind": "search", "pattern": "{alias} S{season:00}E{episode:00}"},
			{"kind": "required", "pattern": "multisub"},
			{"kind": "rejected", "pattern": "hdtv"},
			{"kind": "preferred", "pattern": "web", "score": 30}
		]
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/filter-profile", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.CreateFilterProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	profiles, err := repo.List()
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	assert.Equal(t, "HD", profiles[0].Name)

	_, patterns, err := repo.GetByID(profiles[0].ID)
	require.NoError(t, err)
	require.Len(t, patterns, 4)
	assert.Equal(t, filterProfileModel.KindSearch, patterns[0].Kind)
	assert.Equal(t, 30, patterns[3].Score)
}

func TestCreateFilterProfile_InvalidKind(t *testing.T) {
	h, _ := newFilterProfileTestHandler()

	requestJSON := `{"name": "HD", "patterns": [{"kind": "bogus", "pattern": "x"}]}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/filter-profile", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.CreateFilterProfile(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestCreateFilterProfile_EmptyName(t *testing.T) {
	h, _ := newFilterProfileTestHandler()

	requestJSON := `{"name": "  ", "patterns": []}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/filter-profile", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.CreateFilterProfile(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestUpdateFilterProfile_Success(t *testing.T) {
	h, repo := newFilterProfileTestHandler()

	createdID, err := repo.Create(filterProfileModel.FilterProfile{Name: "Old"}, nil)
	require.NoError(t, err)

	requestJSON := `{"name": "New", "patterns": [{"kind": "search", "pattern": "{alias}"}]}`
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/database/filter-profile/%d", createdID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(createdID, 10))

	err = h.UpdateFilterProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	profile, patterns, err := repo.GetByID(createdID)
	require.NoError(t, err)
	assert.Equal(t, "New", profile.Name)
	require.Len(t, patterns, 1)
}

func TestUpdateFilterProfile_NotFound(t *testing.T) {
	h, _ := newFilterProfileTestHandler()

	requestJSON := `{"name": "New", "patterns": []}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/database/filter-profile/9999", bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	err := h.UpdateFilterProfile(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestDeleteFilterProfile_Success(t *testing.T) {
	h, repo := newFilterProfileTestHandler()

	createdID, err := repo.Create(filterProfileModel.FilterProfile{Name: "HD"}, nil)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/filter-profile/%d", createdID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(createdID, 10))

	err = h.DeleteFilterProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	profiles, err := repo.List()
	require.NoError(t, err)
	assert.Empty(t, profiles)
}

func TestDeleteFilterProfile_NotFound(t *testing.T) {
	h, _ := newFilterProfileTestHandler()

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/database/filter-profile/9999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	err := h.DeleteFilterProfile(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestGetFilterProfile_Success(t *testing.T) {
	h, repo := newFilterProfileTestHandler()

	createdID, err := repo.Create(filterProfileModel.FilterProfile{Name: "HD"}, []filterProfileModel.FilterPattern{
		{Kind: filterProfileModel.KindSearch, Pattern: "{alias} 1080p"},
	})
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/filter-profile/%d", createdID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(createdID, 10))

	err = h.GetFilterProfile(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "{alias} 1080p")
	assert.Contains(t, rec.Body.String(), "HD")
}

func TestGetFilterProfile_NotFound(t *testing.T) {
	h, _ := newFilterProfileTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/filter-profile/9999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("9999")

	err := h.GetFilterProfile(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}

func TestListFilterProfiles_Success(t *testing.T) {
	h, repo := newFilterProfileTestHandler()

	_, err := repo.Create(filterProfileModel.FilterProfile{Name: "HD"}, nil)
	require.NoError(t, err)
	_, err = repo.Create(filterProfileModel.FilterProfile{Name: "SD"}, nil)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/filter-profile", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.ListFilterProfiles(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "HD")
	assert.Contains(t, rec.Body.String(), "SD")
}
