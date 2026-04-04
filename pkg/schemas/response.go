package schemas

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func SendSuccess(c echo.Context, handler string, data any) {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.JSONPretty(200, gin.H{
		"message": fmt.Sprintf("Operation from handler %s successful", handler),
		"data":    data,
	}, "  ")
}

func SendError(c echo.Context, code int, msg string, data ...any) {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.Response().Status = code
	c.JSONPretty(code, gin.H{
		"message": msg,
		"data":    data,
	}, "  ")
}

type ErrorResponse struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

type DefaultResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
