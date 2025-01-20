package controller

import (
	"gateway-service/dto"
	"gateway-service/src/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type paymentController struct {
	paymentService service.PaymentService
}

func NewPaymentController(paymentService service.PaymentService) *paymentController {
	return &paymentController{paymentService}
}

func (h *paymentController) GetWalletByUserId(c echo.Context) error {

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	status, response := h.paymentService.GetWalletByUserId(token)

	return c.JSON(status, response)
}

// CreatePayment godoc
// @Summary Register a new payment
// @Tags payments
// @Accept json
// @Produce json
// @Param order body dto.RegisterRequest true "Register input"
// @Success 201 {object} dto.Payment
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/register [post]
func (h *paymentController) Withdraw(c echo.Context) error {
	var req dto.WithdrawRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.paymentService.Withdraw(token, req)

	return c.JSON(status, response)
}

// LoginPayment godoc
// @Summary Login
// @Tags payments
// @Accept json
// @Produce json
// @Param order body dto.LoginRequest true "Login Information"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/login [post]
func (h *paymentController) CreateInvoice(c echo.Context) error {
	var req dto.CreateInvoiceRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.paymentService.CreateInvoice(token, req)

	return c.JSON(status, response)
}

// LoginPayment godoc
// @Summary Login
// @Tags payments
// @Accept json
// @Produce json
// @Param order body dto.LoginRequest true "Login Information"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/login [post]
func (h *paymentController) XenditInvoiceCallback(c echo.Context) error {
	var req dto.XenditCallback

	callbackToken := c.Request().Header.Get("x-callback-token")
	if callbackToken == "" {
		return c.JSON(http.StatusUnauthorized, "Missing callback token")
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	balanceUpdateReq := dto.UpdateWalletBalanceRequest{
		Amount: req.Amount,
		Type:   "money_in",
	}
	status, response := h.paymentService.UpdateWalletBalance(callbackToken, balanceUpdateReq)

	return c.JSON(status, response)
}

// LoginPayment godoc
// @Summary Login
// @Tags payments
// @Accept json
// @Produce json
// @Param order body dto.LoginRequest true "Login Information"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /payments/login [post]
func (h *paymentController) XenditDisbursementCallback(c echo.Context) error {
	var req dto.XenditCallback

	callbackToken := c.Request().Header.Get("x-callback-token")
	if callbackToken == "" {
		return c.JSON(http.StatusUnauthorized, "Missing callback token")
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	balanceUpdateReq := dto.UpdateWalletBalanceRequest{
		Amount: req.Amount,
		Type:   "money_out",
	}
	status, response := h.paymentService.UpdateWalletBalance(callbackToken, balanceUpdateReq)

	return c.JSON(status, response)
}
