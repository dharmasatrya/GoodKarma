package service

import (
	"context"

	entity "github.com/dharmasatrya/goodkarma/user-service/entity"
	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
)

type UserService interface {
	RegisterUserSupporter(entity.CreateUserSupporterRequest) error
	RegisterUserCoordinator(entity.CreateUserCoordinatorRequest) error
	Login(entity.LoginRequest) (*pb.LoginResponse, error)
	GetUserById(string) (*pb.GetUserByIdResponse, error)
	VerifyEmail(string) (*pb.Empty, error)
}

type userService struct {
	Client pb.UserServiceClient
}

func NewUserService(userClient pb.UserServiceClient) *userService {
	return &userService{userClient}
}

func (us *userService) RegisterUserSupporter(payload entity.CreateUserSupporterRequest) error {
	_, err := us.Client.CreateUserSupporter(context.Background(), &pb.CreateUserSupporterRequest{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		Role:     payload.Role,
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

func (us *userService) RegisterUserCoordinator(payload entity.CreateUserCoordinatorRequest) error {
	_, err := us.Client.CreateUserCoordinator(context.Background(), &pb.CreateUserCoordinatorRequest{
		Username:          payload.Username,
		Email:             payload.Email,
		Password:          payload.Password,
		Role:              payload.Role,
		FullName:          payload.FullName,
		Address:           payload.Address,
		Phone:             payload.Phone,
		Photo:             payload.Photo,
		AccountHolderName: payload.AccountHolderName,
		BankCode:          payload.BankCode,
		BankAccountNumber: payload.BankAccountNumber,
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

func (us *userService) GetUserById(id string) (*pb.GetUserByIdResponse, error) {
	res, err := us.Client.GetUserById(context.Background(), &pb.GetUserByIdRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (us *userService) VerifyEmail(token string) (*pb.Empty, error) {
	_, err := us.Client.VerifyEmail(context.Background(), &pb.VerifyEmailRequest{
		Token: token,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}
