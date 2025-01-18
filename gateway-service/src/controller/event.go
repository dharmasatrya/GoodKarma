package controller

import (
	"gateway-service/dto"
	"gateway-service/src/service"
	"net/http"

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

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.CreateEvent(req)

	return c.JSON(status, response)
}

func (h *eventController) EditEvent(c echo.Context) error {
	var req dto.EventRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.EditEvent(req)

	return c.JSON(status, response)
}

func (h *eventController) GetAllEvents(c echo.Context) error {
	var req dto.EventRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetAllEvent()

	return c.JSON(status, response)
}

func (h *eventController) GetEventById(c echo.Context) error {
	var req dto.EventRequest

	id := c.Param("id")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetEventById(id)

	return c.JSON(status, response)
}

func (h *eventController) GetEventByUserLogin(c echo.Context) error {
	var req dto.EventRequest

	user_id := c.Param("id")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetEventByUserId(user_id)

	return c.JSON(status, response)
}

func (h *eventController) GetEventByCategory(c echo.Context) error {
	var req dto.EventRequest

	id := c.Param("category")
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	status, response := h.eventService.GetEventByCategory(category)

	return c.JSON(status, response)
}
