package service

import (
	"context"
	"gateway-service/dto"
	"log"
	"net/http"

	pb "gateway-service/pb/user"
)

type UserService interface {
	RegisterUser(input dto.RegisterRequest) (int, *dto.User)
	LoginUser(input dto.LoginRequest) (int, *dto.LoginResponse)
}

type userService struct {
	Client pb.UserServiceClient
}

func NewUserService(userClient pb.UserServiceClient) *userService {
	return &userService{userClient}
}

func (u *userService) RegisterUser(input dto.RegisterRequest) (int, *dto.User) {
	res, err := u.Client.RegisterUser(context.Background(), &pb.RegistrationRequest{})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.User{
		Username: res.Username,
	}

	return http.StatusOK, &response
}

func (u *userService) LoginUser(input dto.LoginRequest) (int, *dto.LoginResponse) {
	res, err := u.Client.LoginUser(context.Background(), &pb.LoginRequest{Username: input.Username, Password: input.Password})
	if err != nil {
		log.Fatalf("error while create request %v", err)
	}

	response := dto.LoginResponse{
		Token:        res.Token,
		Success:      true,
		ErrorMessage: "",
	}

	return http.StatusOK, &response
}
