package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	paymentPb "github.com/dharmasatrya/goodkarma/payment-service/proto"

	"github.com/dharmasatrya/goodkarma/user-service/entity"
	"github.com/dharmasatrya/goodkarma/user-service/helper"
	pb "github.com/dharmasatrya/goodkarma/user-service/proto"
	"github.com/dharmasatrya/goodkarma/user-service/repository"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	userRepository repository.UserRepository
	messageBroker  MessageBroker
	paymentClient  paymentPb.PaymentServiceClient
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepository repository.UserRepository, messageBroker MessageBroker, paymentClient paymentPb.PaymentServiceClient) *UserService {
	return &UserService{
		userRepository: userRepository,
		messageBroker:  messageBroker,
		paymentClient:  paymentClient,
	}
}

func (us *UserService) CreateUserSupporter(ctx context.Context, req *pb.CreateUserSupporterRequest) (*pb.CreateUserSupporterResponse, error) {
	payload := entity.CreateUserSupporterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		FullName: req.FullName,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	// Validate the request
	if err := us.validateCreateUserRequest(payload); err != nil {
		return nil, err
	}

	// Create the user
	result, err := us.userRepository.CreateUserSupporter(payload)

	if err != nil {
		return nil, err
	}

	_, err = us.paymentClient.CreateKarma(context.Background(), &paymentPb.CreateKarmaRequest{
		UserId: result.ID.Hex(),
		Amount: 0,
	})

	if req.ReferralCode != "" {
		// Get referral count
		referralCount, err := us.paymentClient.GetReferralCount(context.Background(), &paymentPb.GetReferralCountRequest{
			ReferralCode: req.ReferralCode,
		})

		if err != nil {
			return nil, err
		}

		// Calculate karma
		karmaAmount := uint32(0)

		if referralCount.Count >= 0 {
			karmaAmount = uint32(5000)
		} else if referralCount.Count > 10 && referralCount.Count < 20 {
			karmaAmount = uint32(10000)
		} else {
			karmaAmount = uint32(15000)
		}

		userIdReferrer, err := us.userRepository.GetUserByReferralCode(req.ReferralCode)

		if err != nil {
			return nil, err
		}

		log.Printf("referral count: %v", referralCount.Count)
		log.Printf("karma amount: %v", karmaAmount)
		// Update karma amount for referrer
		_, err = us.paymentClient.UpdateKarmaAmount(context.Background(), &paymentPb.UpdateKarmaAmountRequest{
			UserId: userIdReferrer,
			Amount: karmaAmount,
		})

		if err != nil {
			return nil, err
		}

		// Update karma amount for referee
		_, err = us.paymentClient.UpdateKarmaAmount(context.Background(), &paymentPb.UpdateKarmaAmountRequest{
			UserId: result.ID.Hex(),
			Amount: karmaAmount,
		})

		if err != nil {
			return nil, err
		}

		// Create referral log
		_, err = us.paymentClient.CreateReferralLog(context.Background(), &paymentPb.CreateReferralLogRequest{
			UserId:       result.ID.Hex(),
			ReferralCode: req.ReferralCode,
		})

		if err != nil {
			return nil, err
		}
	}

	tokenString, err := us.generateJWTToken(result)

	// Send email verification
	err = us.sendEmail(req.Email, tokenString)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserSupporterResponse{
		Id: result.ID.Hex(),
	}, nil
}

func (us *UserService) CreateUserCoordinator(ctx context.Context, req *pb.CreateUserCoordinatorRequest) (*pb.CreateUserCoordinatorResponse, error) {
	payload := entity.CreateUserCoordinatorRequest{
		Username:          req.Username,
		Email:             req.Email,
		Password:          req.Password,
		Role:              req.Role,
		FullName:          req.FullName,
		Address:           req.Address,
		Phone:             req.Phone,
		Photo:             req.Photo,
		AccountHolderName: req.AccountHolderName,
		BankCode:          req.BankCode,
		BankAccountNumber: req.BankAccountNumber,
	}

	reqUser := entity.CreateUserSupporterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		FullName: req.FullName,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	reqBank := entity.CreateUserCoordinatorRequest{
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

	// Create the user
	result, err := us.userRepository.CreateUserCoordinator(payload)

	if err != nil {
		return nil, err
	}

	// Create wallet
	err = us.createWallet(result.ID.Hex(), reqBank)

	if err != nil {
		return nil, err
	}

	tokenString, err := us.generateJWTToken(result)

	// Send email verification
	err = us.sendEmail(req.Email, tokenString)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
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

	tokenString, err := us.generateJWTToken(result)

	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Id:    result.ID.Hex(),
		Token: tokenString,
	}, nil
}

func (us *UserService) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {
	result, err := us.userRepository.GetUserById(req.Id)

	if err != nil {
		return nil, err
	}

	return &pb.GetUserByIdResponse{
		Id:       result.ID.Hex(),
		Username: result.Username,
		Email:    result.Email,
		Role:     result.Role,
		FullName: result.FullName,
		Address:  result.Address,
		Phone:    result.Phone,
		Photo:    result.Photo,
	}, nil
}

func (us *UserService) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	payload := entity.UpdateProfileRequest{
		UserID:   req.Id,
		FullName: req.FullName,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	res, err := us.userRepository.UpdateProfile(payload)

	if err != nil {
		return nil, err
	}

	return &pb.UpdateProfileResponse{Id: res.ID.Hex()}, nil
}

func (us *UserService) validateCreateUserRequest(req entity.CreateUserSupporterRequest) error {
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

	if req.FullName == "" {
		return fmt.Errorf("full name is required")
	}

	if req.Address == "" {
		return fmt.Errorf("address is required")
	}

	if req.Phone == "" {
		return fmt.Errorf("phone is required")
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(req.Phone) {
		return fmt.Errorf("phone must contain only numeric characters")
	}

	if len(req.Phone) < 10 || len(req.Phone) > 18 {
		return fmt.Errorf("phone is invalid")
	}

	if !strings.HasPrefix(req.Phone, "0") && !strings.HasPrefix(req.Phone, "62") {
		return fmt.Errorf("phone is invalid")
	}

	return nil
}

func (us *UserService) validateBankRequest(req entity.CreateUserCoordinatorRequest) error {
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

func (us *UserService) generateJWTToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	})

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Println(err)
		return "", err
	}

	return tokenString, nil
}

func (us *UserService) createWallet(userID string, request entity.CreateUserCoordinatorRequest) error {
	// Make the gRPC call to create a wallet
	_, err := us.paymentClient.CreateWallet(context.Background(), &paymentPb.CreateWalletRequest{
		UserId:            userID,
		BankAccountName:   request.AccountHolderName,
		BankCode:          request.BankCode,
		BankAccountNumber: request.BankAccountNumber,
	})

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (us *UserService) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.Empty, error) {
	claims, err := helper.GetClaims(req.Token)

	if err != nil {
		return nil, err
	}

	userID := claims["user_id"].(string)

	err = us.userRepository.VerifyEmail(userID)

	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func (us *UserService) sendEmail(email, tokenString string) error {
	baseUrl := os.Getenv("BASE_URL")

	link := fmt.Sprintf("%v/users/email/verify/%v", baseUrl, tokenString)

	dataJsonRequest := entity.UserRegistData{
		Email: email,
		Link:  link,
	}

	dataJson, err := json.Marshal(dataJsonRequest)
	if err != nil {
		return err
	}

	if err := us.messageBroker.PublishRegistMessage(dataJson); err != nil {
		return err
	}

	return nil
}
