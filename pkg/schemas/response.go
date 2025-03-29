package schemas

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SendSucess(c *gin.Context, handler string, data interface{}) {
	c.Header("Content-type", "application/json")
	c.JSON(200, gin.H{
		"message": fmt.Sprintf("Operation from handler %s successfull", handler),
		"data":    data,
	})
}

func SendError(c *gin.Context, code int, msg string) {
	c.Header("Content-type", "application/json")
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
