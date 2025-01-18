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

func (us *UserService) CreateUserSupporter(ctx context.Context, req *pb.CreateUserSupporterRequest) (*pb.CreateUserSupporterResponse, error) {
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

	if err := us.validateCreateUserRequest(payload); err != nil {
		return nil, err
	}

	result, err := us.userRepository.CreateUser(payload)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserSupporterResponse{
		Id: result.ID.Hex(),
	}, nil
}

func (us *UserService) CreateUserCoordinator(ctx context.Context, req *pb.CreateUserCoordinatorRequest) (*pb.CreateUserCoordinatorResponse, error) {
	reqUser := entity.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		FullName: req.FullName,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	// bankName, err := us.getBankName(req.BankCode)
	// if err != nil {
	// 	return nil, err
	// }

	reqBank := entity.CreateMerchantRequest{
		AccountHolderName: req.AccountHolderName,
		BankCode:          req.BankCode,
		BankAccountNumber: req.BankAccountNumber,
	}

	if err := us.validateCreateUserRequest(reqUser); err != nil {
		return nil, err
	}

	if err := us.validateBankRequest(reqBank); err != nil {
		return nil, err
	}

	result, err := us.userRepository.CreateMerchant(reqUser, reqBank)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserCoordinatorResponse{
		Id: result.ID.Hex(),
	}, nil
}

func (us *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	payload := entity.LoginRequest{
		UsernameOrEmail: req.UsernameOrEmail,
		Password:        req.Password,
	}

	result, err := us.userRepository.Login(payload)

	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Id: result.ID.Hex(),
	}, nil
}

func (us *UserService) validateCreateUserRequest(req entity.CreateUserRequest) error {
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

func (us *UserService) validateBankRequest(req entity.CreateMerchantRequest) error {
	if req.AccountHolderName == "" {
		return fmt.Errorf("account holder name is required")
	}

	if req.BankCode == "" {
		return fmt.Errorf("bank code is required")
	}

	if req.BankAccountNumber == "" {
		return fmt.Errorf("bank account number is required")
	}

	return nil
}

// func (us *UserService) getBankName(bankCode string) (string, error) {
// 	bankList, err := helper.GetBankList()
// 	if err != nil {
// 		return "", err
// 	}

// 	for _, bank := range *bankList {
// 		if bank.Code == strings.ToUpper(bankCode) {
// 			return bank.Name, nil
// 		}
// 	}

// 	return "", fmt.Errorf("bank code is invalid")
// }
