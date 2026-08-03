package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	showSchema "github.com/jusoaresg/gorgon/internal/show/schema"
	"github.com/jusoaresg/gorgon/testutils"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddShowToList_Success(t *testing.T) {
	e := echo.New()
	_ = e

	request := showSchema.AddShowToListRequest{
		Id:           10,
		TrackingType: "all",
	}
	requestJSON, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/show", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	logger := slog.Default()
	db := testutils.GetTestDB()

	h := &Handler{
		Logger: logger,
		DB:     db,
	}

	_, err = h.addShowToListHandler(c, &request)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
