package controller

import (
	"gateway-service/dto"
	"gateway-service/helpers"
	"gateway-service/src/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type eventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) *eventController {
	return &eventController{eventService}
}

// CreateEvent godoc
// @Summary Create a new event
// @Tags events
// @Accept json
// @Produce json
// @Param order body dto.CreateEventRequest true "Event input"
// @Success 201 {object} dto.Event
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /events [post]
func (h *eventController) CreateEvent(c echo.Context) error {
	var req dto.EventRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response, err := h.eventService.CreateEvent(token, req)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	if status != http.StatusCreated {
		responseError := helpers.AssertJSONStatus(status)
		return c.JSON(status, responseError)
	}

	return c.JSON(status, response)
}

func (h *eventController) EditEvent(c echo.Context) error {
	var req dto.UpdateDescriptionRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	id := c.Param("id")
	idConv, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id params")
	}

	status, response := h.eventService.EditEvent(token, idConv, req)
	if status != http.StatusOK {
		responseError := helpers.AssertJSONStatus(status)
		return c.JSON(status, responseError)
	}

	return c.JSON(status, response)
}

func (h *eventController) GetAllEvents(c echo.Context) error {
	var req dto.EventRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetAllEvents()

	return c.JSON(status, response)
}

func (h *eventController) GetEventById(c echo.Context) error {
	var req dto.EventRequest

	id := c.Param("id")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	idConv, err := strconv.Atoi(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid id params")
	}

	status, response := h.eventService.GetEventById(idConv)
	if status != http.StatusOK {
		responseError := helpers.AssertJSONStatus(status)
		return c.JSON(status, responseError)
	}

	return c.JSON(status, response)
}

func (h *eventController) GetAllEventByUserLogin(c echo.Context) error {
	var req dto.EventRequest

	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "No token provided",
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetEventByUserLogin(token)
	if status != http.StatusOK {
		responseError := helpers.AssertJSONStatus(status)
		return c.JSON(status, responseError)
	}

	return c.JSON(status, response)
}

func (h *eventController) GetAllEventByCategory(c echo.Context) error {
	var req dto.EventRequest

	category := c.QueryParam("category")

	// Jika kategori tidak diberikan, alihkan ke handler lain atau kirimkan error
	if category == "" {
		return c.JSON(http.StatusBadRequest, "Invalid category param") // Memanggil rute lainnya yang mengambil semua events
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetEventByCategory(category)
	if status != http.StatusOK {
		responseError := helpers.AssertJSONStatus(status)
		return c.JSON(status, responseError)
	}

	return c.JSON(status, response)
}
