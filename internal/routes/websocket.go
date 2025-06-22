package routes

import (
	"github.com/jusoaresg/gorgon/config"
	"github.com/jusoaresg/gorgon/pkg/handler"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func SetupWebsocketRouter(r *echo.Group) {
	logger := config.GetLogger()

	r.GET("ws", handler.WebSocketHandler)
	logger.Info("Websocket router added successfully", slog.String("endpoint", "/api/v1/ws"))
}
