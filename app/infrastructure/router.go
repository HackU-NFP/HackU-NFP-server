package infrastructure

import (
	"nfp-server/interfaces/controllers"

	"github.com/labstack/echo/v4"
)

// Router ルーティング
type Router struct {
	e  *echo.Echo
	lc *controllers.LinebotController
	ac *controllers.ApiController
}

// NewRouter コンストラクタ
func NewRouter(e *echo.Echo, lc *controllers.LinebotController, ac *controllers.ApiController) *Router {
	return &Router{e: e, lc: lc, ac: ac}
}

func (r *Router) Init() {
	r.e.POST("/linebot/callback", r.lc.CatchEvents())
	r.e.GET("/api/nfts", r.ac.GetNfts())
	r.e.GET("/api/nft", r.ac.GetNft())
}
