package controller

import (
	"gateway-service/dto"
	"gateway-service/src/service"

	"github.com/labstack/echo/v4"
)

type KarmaController struct {
	karmaService service.KarmaService
}

func NewKarmaController(karmaService service.KarmaService) *KarmaController {
	return &KarmaController{karmaService}
}

func (kc *KarmaController) GetKarmaReward(c echo.Context) error {
	res, err := kc.karmaService.GetKarmaReward()

	if err != nil {
		return c.JSON(500, dto.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, dto.ResponseOK{
		Message: "Success",
		Data:    res,
	})
}

func (kc *KarmaController) ExchangeReward(c echo.Context) error {
	id := c.Param("id")
	jwtToken := c.Request().Header.Get("Authorization")

	if jwtToken == "" {
		return c.JSON(401, dto.ResponseError{
			Message: "No token provided",
		})
	}

	err := kc.karmaService.ExchangeReward(jwtToken, id)

	if err != nil {
		return c.JSON(500, dto.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(200, dto.ResponseOK{
		Message: "Success",
		Data:    nil,
	})
}
