package infrastructure

import (
	"net/http"
	"nfp-server/interfaces/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() {
	e := echo.New()
	// ミドルウェアを設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	linebotController := controllers.NewLinebotController()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/linebot/callback", linebotController.CatchEvents())
	e.Logger.Fatal(e.Start(":8080"))
}
