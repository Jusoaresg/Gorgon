package schemas

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

func SendSuccess(c echo.Context, handler string, data interface{}) {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.JSONPretty(200, gin.H{
		"message": fmt.Sprintf("Operation from handler %s successful", handler),
		"data":    data,
	}, "  ")
}

func SendError(c echo.Context, code int, msg string) {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	c.JSON(code, gin.H{
		"message":   msg,
		"errorCode": code,
	})
}

type ErrorResponse struct {
	Message   string `json:"message"`
	ErrorCode string `json:"errorCode"`
}

type DefaultResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
