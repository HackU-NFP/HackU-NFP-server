//go:build wireinject
// +build wireinject

package main

import (
	"nfp-server/infrastructure"
	"nfp-server/interfaces/controllers"
	"nfp-server/interfaces/presenter"
	"nfp-server/usecase"
	"nfp-server/usecase/ipresenter"

	"github.com/google/wire"
	"github.com/labstack/echo/v4"
)

var superSet = wire.NewSet(
	// Presenter
	presenter.NewLinePresenter,
	wire.Bind(new(ipresenter.ILinePresenter), new(*presenter.LinePresenter)),

	// Interactor
	usecase.NewBlockchainInteractor,
	usecase.NewLineBotInteractor,
	wire.Bind(new(usecase.IBlockchainUseCase), new(*usecase.BlockchainInteractor)),
	wire.Bind(new(usecase.ILineBotUseCase), new(*usecase.LineBotInteractor)),

	// Controller
	controllers.NewLinebotController,

	// Router
	infrastructure.NewRouter,
)

// Initialize DI
func Initialize(e *echo.Echo) *infrastructure.Router {
	wire.Build(superSet)
	return &infrastructure.Router{}
}
