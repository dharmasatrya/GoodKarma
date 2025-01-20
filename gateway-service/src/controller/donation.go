package controller

import (
	"fmt"
	"gateway-service/dto"
	"gateway-service/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type donationController struct {
	donationService service.DonationService
}

func NewDonationController(donationService service.DonationService) *donationController {
	return &donationController{donationService}
}

// CreateDonation godoc
// @Summary Create a new donation
// @Tags donations
// @Accept json
// @Produce json
// @Param order body dto.CreateDonationRequest true "Donation input"
// @Success 201 {object} dto.Donation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /donations [post]
func (h *donationController) CreateDonation(c echo.Context) error {
	var req dto.CreateDonationRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.donationService.CreateDonation(token, req)

	return c.JSON(status, response)
}

// LoginDonation godoc
// @Summary Login
// @Tags donations
// @Accept json
// @Produce json
// @Param order body dto.LoginRequest true "Login Information"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /donations/:id [put]
// benerin docs ntaran
func (h *donationController) UpdateDonationStatus(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	var req dto.UpdateDonationStatusRequest

	req.ID = c.Param("id")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.donationService.UpdateDonationStatus(token, req)

	return c.JSON(status, response)
}

func (h *donationController) GetAllDonationByUser(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	status, response := h.donationService.GetAllDonationByUser(token)

	return c.JSON(status, response)
}

func (h *donationController) GetAllDonationByEventId(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	eventID := c.Param("event_id")

	status, response := h.donationService.GetAllDonationByEventId(token, eventID)

	return c.JSON(status, response)
}
