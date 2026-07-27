package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jusoaresg/gorgon/internal/indexer/model"
	"github.com/jusoaresg/gorgon/internal/indexer/repository"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newIndexerTestHandler() (*Handler, *repository.IndexerRepository) {
	db := testutils.GetTestDB()
	indexerRepo := repository.NewIndexerRepository(db)
	h := &Handler{
		IndexerRepo: indexerRepo,
		Logger:      slog.Default(),
	}
	return h, indexerRepo
}

func createTestIndexer(t *testing.T, indexerRepo *repository.IndexerRepository) int64 {
	t.Helper()
	indexer := model.Indexer{
		IndexerID:      1,
		Name:           "Test Indexer",
		DefinitionName: "test",
		IndexerUrls:    "[]",
		Language:       "en",
	}
	err := indexerRepo.Create(indexer)
	assert.NoError(t, err)

	created, err := indexerRepo.GetById(1)
	assert.NoError(t, err)
	return created.ID
}

func TestListIndexers_Success(t *testing.T) {
	h, indexerRepo := newIndexerTestHandler()

	createTestIndexer(t, indexerRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/indexer", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.ListIndexers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestListIndexers_Empty(t *testing.T) {
	h, _ := newIndexerTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/indexer", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.ListIndexers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
	assert.Contains(t, rec.Body.String(), "List Indexers")
}

func TestGetIndexer_Success(t *testing.T) {
	h, indexerRepo := newIndexerTestHandler()

	createTestIndexer(t, indexerRepo)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/indexer/1", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	err := h.GetIndexer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestGetIndexer_NotFound(t *testing.T) {
	h, _ := newIndexerTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/database/indexer/999999", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("999999")

	err := h.GetIndexer(c)
	assert.Error(t, err)
}

func TestDeleteIndexer_Success(t *testing.T) {
	h, indexerRepo := newIndexerTestHandler()

	createTestIndexer(t, indexerRepo)

	request := struct {
		Id int64 `json:"id"`
	}{Id: 1}
	requestJSON, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/database/indexer", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.DeleteIndexer(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestDeleteIndexer_NotFound(t *testing.T) {
	h, _ := newIndexerTestHandler()

	request := struct {
		Id int64 `json:"id"`
	}{Id: 999999}
	requestJSON, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/database/indexer", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.DeleteIndexer(c)
	assert.Error(t, err)
}
