package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jusoaresg/gorgon/internal/show/repository"
	"github.com/jusoaresg/gorgon/pkg/schemas"
	"github.com/jusoaresg/gorgon/testutils"


	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newShowTestHandler() (*Handler, *repository.ShowRepository) {
	db := testutils.GetTestDB()
	showRepo := repository.NewShowRepository(db)
	h := &Handler{
		ShowRepo: showRepo,
		Logger:   slog.Default(),
		DB:       db,
	}
	return h, showRepo
}

func TestGetShow_Success(t *testing.T) {
	h, showRepo := newShowTestHandler()

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/database/show/%d", id), nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(id, 10))

	err = h.GetShow(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), show.Name)
}

func TestGetShow_InvalidId(t *testing.T) {
	h, _ := newShowTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/abc", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := h.GetShow(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestGetShow_NotFound(t *testing.T) {
	h, _ := newShowTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show/999999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999999")

	err := h.GetShow(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Result().StatusCode)
}

func TestListShows_Success(t *testing.T) {
	h, showRepo := newShowTestHandler()

	show1 := testutils.MakeFakeShow()
	show2 := testutils.MakeFakeShow()
	_, err := showRepo.Create(show1)
	assert.NoError(t, err)
	_, err = showRepo.Create(show2)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err = h.ListShows(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestListShows_Empty(t *testing.T) {
	h, _ := newShowTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/show", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.ListShows(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "List Shows")
}

func TestDeleteShow_InvalidId(t *testing.T) {
	h, _ := newShowTestHandler()

	request := schemas.IdRequest{Id: -100}
	requestJSON, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/database/show/-100", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("-100")

	err := h.DeleteShow(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Result().StatusCode)
}

func TestDeleteShow_HTMX(t *testing.T) {
	h, showRepo := newShowTestHandler()

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	request := schemas.IdRequest{Id: id}
	requestJSON, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d", id), bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("HX-Request", "true")
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(id, 10))

	err = h.DeleteShow(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Equal(t, "/", rec.Header().Get("HX-Redirect"))
}

func TestAddShowToList_InvalidTracking(t *testing.T) {
	db := testutils.GetTestDB()
	h := &Handler{
		Logger: slog.Default(),
		DB:     db,
	}

	request := schemas.IdRequest{Id: 10}
	requestJSON, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/show", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.AddShowToList(c)
	assert.Error(t, err)
}
