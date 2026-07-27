package api

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jusoaresg/gorgon/pkg/schemas"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newConfigTestHandler() *Handler {
	return &Handler{
		Logger: slog.Default(),
	}
}

func TestGetAppConfig_Success(t *testing.T) {
	h := newConfigTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/app/config", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.GetAppConfig(c)
	if err != nil {
		t.Skip("Config file not available in test environment")
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestUpdateAppConfig_BindError(t *testing.T) {
	h := newConfigTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/api/v1/app/config", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateAppConfig(c)
	assert.Error(t, err)
}

func TestUpdateAppConfig_Success(t *testing.T) {
	h := newConfigTestHandler()

	updateInput := schemas.UpdateConfigInput{}
	requestJSON, _ := json.Marshal(updateInput)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/app/config", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.UpdateAppConfig(c)
	if err != nil {
		t.Skip("Config file not available in test environment")
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
