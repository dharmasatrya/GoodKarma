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
