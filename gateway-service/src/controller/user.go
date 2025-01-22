package controller

import (
	"gateway-service/src/service"
	"net/http"

	entity "github.com/dharmasatrya/goodkarma/user-service/entity"

	"github.com/labstack/echo/v4"
)

type userController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *userController {
	return &userController{userService}
}

func (us *userController) RegisterUserSupporter(c echo.Context) error {
	var payload entity.CreateUserSupporterRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Message: "Invalid request payload",
		})
	}

	err := us.userService.RegisterUserSupporter(payload)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, entity.ResponseOK{
		Message: "User has been created",
		Data:    nil,
	})
}

func (us *userController) RegisterUserCoordinator(c echo.Context) error {
	var payload entity.CreateUserCoordinatorRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Message: "Invalid request payload",
		})
	}

	err := us.userService.RegisterUserCoordinator(payload)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, entity.ResponseOK{
		Message: "User has been created",
		Data:    nil,
	})
}

func (us *userController) Login(c echo.Context) error {
	var payload entity.LoginRequest

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, entity.ResponseError{
			Message: "Invalid request payload",
		})
	}

	result, err := us.userService.Login(payload)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.ResponseOK{
		Message: "Login success",
		Data:    result,
	})
}

func (us *userController) GetUserById(c echo.Context) error {
	id := c.Param("id")

	result, err := us.userService.GetUserById(id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.ResponseOK{
		Message: "User found",
		Data:    result,
	})
}

func (us *userController) VerifyEmail(c echo.Context) error {
	token := c.Param("token")

	_, err := us.userService.VerifyEmail(token)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, entity.ResponseError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, entity.ResponseOK{
		Message: "Email has been verified",
		Data:    nil,
	})
}
