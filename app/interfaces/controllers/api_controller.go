package controllers

import (
	"nfp-server/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// LinebotController LINEBOTコントローラ
type ApiController struct {
	blockchainInteractor usecase.IBlockchainUseCase
}

// NewLinebotController コンストラクタ
func NewApiController(blockchainInteractor usecase.IBlockchainUseCase) *ApiController {
	return &ApiController{
		blockchainInteractor: blockchainInteractor,
	}
}

// NFT一覧取得
func (controller *ApiController) GetNfts() echo.HandlerFunc {
	return func(c echo.Context) error {
		contractId := c.QueryParam("contractId")
		orderBy := c.QueryParam("orderBy")
		limit := c.QueryParam("limit")
		page := c.QueryParam("page")

		response, err := controller.blockchainInteractor.GetNonFungibles(contractId, orderBy, limit, page)
		if err != nil {
			logrus.Fatalf("Error LINEBOT parsing request: %v", err)
			return c.JSON(500, NewError(err))
		}

		return c.JSON(200, response)
	}
}

func (controller *ApiController) GetNft() echo.HandlerFunc {
	return func(c echo.Context) error {
		contractId := c.QueryParam("contractId")
		tokenType := c.QueryParam("tokenType")
		response, err := controller.blockchainInteractor.GetNonFungibleInfo(contractId, tokenType)
		if err != nil {

			logrus.Fatalf("Error LINEBOT parsing request: %v", err)
			return c.JSON(500, NewError(err))
		}

		return c.JSON(200, response)
	}
}
