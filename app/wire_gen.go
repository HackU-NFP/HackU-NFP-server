// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"nfp-server/infrastructure"
	"nfp-server/interfaces/controllers"
	"nfp-server/interfaces/presenter"
	"nfp-server/usecase"
	"nfp-server/usecase/ipresenter"
)

// Injectors from wire.go:

// Initialize DI
func Initialize(e *echo.Echo) *infrastructure.Router {
	linePresenter := presenter.NewLinePresenter()
	lineBotInteractor := usecase.NewLineBotInteractor(linePresenter)
	linebotController := controllers.NewLinebotController(lineBotInteractor)
	router := infrastructure.NewRouter(e, linebotController)
	return router
}

// wire.go:

var superSet = wire.NewSet(presenter.NewLinePresenter, wire.Bind(new(ipresenter.ILinePresenter), new(*presenter.LinePresenter)), usecase.NewLineBotInteractor, wire.Bind(new(usecase.ILineBotUseCase), new(*usecase.LineBotInteractor)), controllers.NewLinebotController, infrastructure.NewRouter)