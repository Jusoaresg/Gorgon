package routes

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/jusoaresg/gorgon/assets"
	"github.com/labstack/echo/v4"
)

func SetupFrontRouter(e *echo.Echo) {
	buildFS, err := fs.Sub(assets.FrontStaticFS, "front/build")
	if err != nil {
		panic(fmt.Errorf("error accessing embedded front files: %w", err))
	}

	e.GET("/favicon.png", echo.WrapHandler(http.FileServer(http.FS(buildFS))))
	e.GET("/_app/*", echo.WrapHandler(http.FileServer(http.FS(buildFS))))
	e.GET("/css/*", echo.WrapHandler(http.FileServer(http.FS(buildFS))))

	e.GET("/*", func(c echo.Context) error {
		index, err := fs.ReadFile(buildFS, "index.html")
		if err != nil {
			panic(err)
		}
		return c.Blob(http.StatusOK, echo.MIMETextHTML, index)
	})
}
