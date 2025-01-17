package service

import (
	"context"
	"goodkarma-user-service/entity"
	pb "goodkarma-user-service/proto"
	"goodkarma-user-service/repository"
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

	result, err := us.userRepository.CreateUser(payload)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Id: result.ID.Hex(),
	}, nil
}
