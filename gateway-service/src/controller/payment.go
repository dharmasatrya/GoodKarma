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

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.paymentService.Withdraw(req)

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

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.paymentService.CreateInvoice(req)

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

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	balanceUpdateReq := dto.UpdateWalletBalanceRequest{
		Amount: req.Amount,
		Type:   "money_in",
	}
	status, response := h.paymentService.UpdateWalletBalance(balanceUpdateReq)

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

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	balanceUpdateReq := dto.UpdateWalletBalanceRequest{
		Amount: req.Amount,
		Type:   "money_out",
	}
	status, response := h.paymentService.UpdateWalletBalance(balanceUpdateReq)

	return c.JSON(status, response)
}
