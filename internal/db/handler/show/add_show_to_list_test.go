package show

import (
	"bytes"
	"encoding/json"
	"github.com/jusoaresg/gorgon/internal/db/schema/show"
	"github.com/jusoaresg/gorgon/testutils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddShowToList_Sucess(t *testing.T) {
	e := echo.New()
	_ = e

	request := show.AddShowToListRequest{
		Id:           10,
		TrackingType: "all",
	}
	requestJSON, err := json.Marshal(request)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/database/show", bytes.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	err = addShowToListHandler(c, testutils.GetTestDB())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Result().StatusCode)
}
