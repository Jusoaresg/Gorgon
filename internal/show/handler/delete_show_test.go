package show

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func TestDeleteShow_Success(t *testing.T) {
	db := testutils.GetTestDB()
	showRepo := repository.NewShowRepository(db)

	show := testutils.MakeFakeShow()
	id, err := showRepo.Create(show)
	assert.NoError(t, err)

	var request schemas.IdRequest
	request.Id = id

	requestJSON, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d", id), bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(id, 10))

	err = deleteShowHandler(c, db)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}

func TestDeleteShow_NotFound(t *testing.T) {
	db := testutils.GetTestDB()

	request := schemas.IdRequest{
		Id: 1,
	}
	requestJSON, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/database/show/%d", request.Id), bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(strconv.FormatInt(request.Id, 10))

	err = deleteShowHandler(c, db)
	assert.Error(t, err)
	assert.Equal(t, 404, rec.Result().StatusCode)
}

func TestDeleteShow_IdGreaterThanZero(t *testing.T) {
	db := testutils.GetTestDB()

	request := schemas.IdRequest{
		Id: -100,
	}

	requestJSON, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/database/show", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	err = deleteShowHandler(c, db)
	assert.Error(t, err)
	assert.Equal(t, 400, rec.Result().StatusCode)
}
