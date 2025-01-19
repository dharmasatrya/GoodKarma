package service

import (
	"context"

	entity "github.com/dharmasatrya/goodkarma/user-service/entity"
	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
)

type UserService interface {
	RegisterUserSupporter(entity.CreateUserRequest) error
	RegisterUserCoordinator(entity.CreateUserRequest, entity.CreateMerchantRequest) error
	Login(entity.LoginRequest) (*pb.LoginResponse, error)
}

type userService struct {
	Client pb.UserServiceClient
}

func NewUserService(userClient pb.UserServiceClient) *userService {
	return &userService{userClient}
}

func (us *userService) RegisterUserSupporter(payload entity.CreateUserRequest) error {
	_, err := us.Client.CreateUserSupporter(context.Background(), &pb.CreateUserSupporterRequest{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     "supporter",
		FullName: payload.FullName,
		Address:  payload.Address,
		Phone:    payload.Phone,
		Photo:    payload.Photo,
	})

	if err != nil {
		return err
	}

	return nil
}

func (us *userService) RegisterUserCoordinator(payload entity.CreateUserRequest, payloadWallet entity.CreateMerchantRequest) error {
	_, err := us.Client.CreateUserCoordinator(context.Background(), &pb.CreateUserCoordinatorRequest{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     "organizer",
		FullName: payload.FullName,
		Address:  payload.Address,
		Phone:    payload.Phone,
		Photo:    payload.Photo,
	})

	if err != nil {
		return err
	}

	return nil
}

func (us *userService) Login(input entity.LoginRequest) (*pb.LoginResponse, error) {
	res, err := us.Client.Login(context.Background(), &pb.LoginRequest{
		UsernameOrEmail: input.UsernameOrEmail,
		Password:        input.Password,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

// func (u *userService) RegisterUser(input dto.RegisterRequest) (int, *dto.User) {
// 	res, err := u.Client.RegisterUser(context.Background(), &pb.RegistrationRequest{})
// 	if err != nil {
// 		log.Fatalf("error while create request %v", err)
// 	}

// 	response := dto.User{
// 		Username: res.Username,
// 	}

// 	return http.StatusOK, &response
// }

// func (u *userService) LoginUser(input dto.LoginRequest) (int, *dto.LoginResponse) {
// 	res, err := u.Client.LoginUser(context.Background(), &pb.LoginRequest{Username: input.Username, Password: input.Password})
// 	if err != nil {
// 		log.Fatalf("error while create request %v", err)
// 	}

// 	response := dto.LoginResponse{
// 		Token:        res.Token,
// 		Success:      true,
// 		ErrorMessage: "",
// 	}

// 	return http.StatusOK, &response
// }
