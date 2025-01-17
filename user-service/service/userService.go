package service

import (
	"context"
	"fmt"
	"goodkarma-user-service/entity"
	pb "goodkarma-user-service/proto"
	"goodkarma-user-service/repository"
	"strings"
)

type UserService struct {
	userRepository repository.UserRepository
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (us *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	payload := entity.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		FullName: req.FullName,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	if err := us.validateCreateUserRequest(req); err != nil {
		return nil, err
	}

	result, err := us.userRepository.CreateUser(payload)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id: result.ID.Hex(),
	}, nil
}

func (us *UserService) validateCreateUserRequest(req *pb.CreateUserRequest) error {
	if req.Username == "" {
		return fmt.Errorf("username is required")
	}

	if len(req.Username) < 5 {
		return fmt.Errorf("username must be at least 5 characters")
	}

	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	if len(req.Email) < 8 || len(req.Email) > 50 {
		return fmt.Errorf("email is invalid")
	}

	if req.Password == "" {
		return fmt.Errorf("password is required")
	}

	if len(req.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	if req.Role == "" {
		return fmt.Errorf("role is required")
	}

	if req.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if req.Address == "" {
		return fmt.Errorf("address is required")
	}

	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}

	if len(req.Phone) < 10 || len(req.Phone) > 18 {
		return fmt.Errorf("phone is invalid")
	}

	if !strings.HasPrefix(req.Phone, "0") && !strings.HasPrefix(req.Phone, "62") {
		return fmt.Errorf("phone is invalid")
	}

	return nil
}
