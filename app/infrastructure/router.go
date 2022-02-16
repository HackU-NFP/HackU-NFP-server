package infrastructure

import (
	"nfp-server/interfaces/controllers"

	"github.com/labstack/echo/v4"
)

// Router ルーティング
type Router struct {
	e  *echo.Echo
	lc *controllers.LinebotController
}

// NewRouter コンストラクタ
func NewRouter(e *echo.Echo, lc *controllers.LinebotController) *Router {
	return &Router{e: e, lc: lc}
}

func (r *Router) Init() {
	r.e.POST("/linebot/callback", r.lc.CatchEvents())
}
