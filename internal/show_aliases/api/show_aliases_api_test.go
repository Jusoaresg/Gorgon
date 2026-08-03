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
	showAliasRepository "github.com/jusoaresg/gorgon/internal/show_aliases/repository"
	"github.com/jusoaresg/gorgon/internal/show_aliases/model"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newShowAliasesTestHandler() (*Handler, showAliasRepository.ShowAliasesRepositoryInterface, *showRepository.ShowRepository) {
	db := testutils.GetTestDB()
	aliasRepo := showAliasRepository.NewShowAliasesRepository(db)
	showRepo := showRepository.NewShowRepository(db)
	h := &Handler{
		ShowAliasRepo: &aliasRepo,
		Logger:        slog.Default(),
	}
	return h, &aliasRepo, showRepo
}

func TestGetShowAliases_Success(t *testing.T) {
	h, aliasRepo, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	alias := model.ShowAlias{
		ShowID:  showID,
		Alias:   "Test Alias",
		Country: "US",
		Source:  "tvmaze",
	}
	_, err = aliasRepo.Create(alias)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show/aliases/%d", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.GetShowAliases(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "Test Alias")
}

func TestGetShowAliases_InvalidId(t *testing.T) {
	h, _, _ := newShowAliasesTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/aliases/abc", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := h.GetShowAliases(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestGetShowAliases_Empty(t *testing.T) {
	h, _, _ := newShowAliasesTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/aliases/999999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999999")

	err := h.GetShowAliases(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "Get Show Aliases")
}

func TestAddShowAlias_Success(t *testing.T) {
	h, aliasRepo, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	requestJSON := `{"alias": "DBZ"}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/database/show/%d/alias", showID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.AddShowAlias(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	aliases, err := aliasRepo.ListByShowID(showID)
	assert.NoError(t, err)
	require.Len(t, aliases, 1)
	assert.Equal(t, "DBZ", aliases[0].Alias)
	assert.Equal(t, "user", aliases[0].Source)
}

func TestAddShowAlias_EmptyAlias(t *testing.T) {
	h, _, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	requestJSON := `{"alias": "   "}`
	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/v1/database/show/%d/alias", showID), bytes.NewReader([]byte(requestJSON)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(showID, 10))

	err = h.AddShowAlias(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestDeleteShowAlias_Success(t *testing.T) {
	h, aliasRepo, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	aliasID, err := aliasRepo.Create(model.ShowAlias{ShowID: showID, Alias: "DBZ", Source: "user"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d/alias/%d", showID, aliasID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "aliasId")
	c.SetParamValues(strconv.FormatInt(showID, 10), strconv.FormatInt(aliasID, 10))

	err = h.DeleteShowAlias(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)

	aliases, err := aliasRepo.ListByShowID(showID)
	assert.NoError(t, err)
	assert.Empty(t, aliases)
}

func TestDeleteShowAlias_OnlyUserAliases(t *testing.T) {
	h, aliasRepo, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	aliasID, err := aliasRepo.Create(model.ShowAlias{ShowID: showID, Alias: "DBZ", Source: "tvmaze"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d/alias/%d", showID, aliasID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "aliasId")
	c.SetParamValues(strconv.FormatInt(showID, 10), strconv.FormatInt(aliasID, 10))

	err = h.DeleteShowAlias(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestDeleteShowAlias_WrongShow(t *testing.T) {
	h, aliasRepo, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)
	otherShowID, err := showRepo.Create(testutils.MakeFakeShow())
	assert.NoError(t, err)

	aliasID, err := aliasRepo.Create(model.ShowAlias{ShowID: showID, Alias: "DBZ", Source: "user"})
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d/alias/%d", otherShowID, aliasID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "aliasId")
	c.SetParamValues(strconv.FormatInt(otherShowID, 10), strconv.FormatInt(aliasID, 10))

	err = h.DeleteShowAlias(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestDeleteShowAlias_NotFound(t *testing.T) {
	h, _, showRepo := newShowAliasesTestHandler()

	show := testutils.MakeFakeShow()
	showID, err := showRepo.Create(show)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d/alias/99999", showID), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id", "aliasId")
	c.SetParamValues(strconv.FormatInt(showID, 10), "99999")

	err = h.DeleteShowAlias(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Result().StatusCode)
}
