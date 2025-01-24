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

	karmaPb "github.com/dharmasatrya/goodkarma/karma-service/proto"
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
	karmaClient    karmaPb.KarmaServiceClient
	pb.UnimplementedUserServiceServer
}

func NewUserService(userRepository repository.UserRepository, messageBroker MessageBroker, paymentClient paymentPb.PaymentServiceClient, karmaClient karmaPb.KarmaServiceClient) *UserService {
	return &UserService{
		userRepository: userRepository,
		messageBroker:  messageBroker,
		paymentClient:  paymentClient,
		karmaClient:    karmaClient,
	}
}

func (us *UserService) CreateUserSupporter(ctx context.Context, req *pb.CreateUserSupporterRequest) (*pb.CreateUserSupporterResponse, error) {
	payload := entity.CreateUserSupporterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     "supporter",
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

	_, err = us.karmaClient.CreateKarma(context.Background(), &karmaPb.CreateKarmaRequest{
		UserId: result.ID.Hex(),
		Amount: 0,
	})

	log.Println("Referral code:", req.ReferralCode)
	// If there is a referral code, process it
	if req.ReferralCode != "" {
		if err := us.ProcessReferral(ctx, req.ReferralCode, result.ID.Hex()); err != nil {
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
		Role:              "coordinator",
		FullName:          req.FullName,
		Address:           req.Address,
		Phone:             req.Phone,
		Photo:             req.Photo,
		NIK:               req.Nik,
		AccountHolderName: req.AccountHolderName,
		BankCode:          req.BankCode,
		BankAccountNumber: req.BankAccountNumber,
	}

	reqUser := entity.CreateUserSupporterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
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

	if err := helper.ValidateNIK(req.Nik); err != nil {
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
	err = us.CreateWallet(result.ID.Hex(), reqBank)

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

func (us *UserService) CreateWallet(userID string, request entity.CreateUserCoordinatorRequest) error {
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

// processReferral handles all referral-related logic.
func (us *UserService) ProcessReferral(ctx context.Context, referralCode, refereeUserID string) error {
	// Get referrer user ID
	referrerUserID, err := us.userRepository.GetUserByReferralCode(referralCode)
	if err != nil {
		return err
	}

	// Get referral count
	referralCount, err := us.karmaClient.GetReferralCount(ctx, &karmaPb.GetReferralCountRequest{
		ReferralCode: referralCode,
	})
	if err != nil {
		return err
	}

	// Calculate karma amount
	karmaAmountReferrer := calculateKarmaAmount(referralCount.Count)
	karmaAmmountReferee := karmaAmountReferrer * 25 / 100

	// Update karma for referrer
	if _, err := us.karmaClient.UpdateKarmaAmount(ctx, &karmaPb.UpdateKarmaAmountRequest{
		UserId: referrerUserID,
		Amount: karmaAmountReferrer,
	}); err != nil {
		return err
	}

	// Update karma for referee
	if _, err := us.karmaClient.UpdateKarmaAmount(ctx, &karmaPb.UpdateKarmaAmountRequest{
		UserId: refereeUserID,
		Amount: karmaAmmountReferee,
	}); err != nil {
		return err
	}

	log.Println("Referee user id:", refereeUserID)
	log.Println("referral code:", referralCode)
	// Create referral log
	if _, err := us.karmaClient.CreateReferralLog(ctx, &karmaPb.CreateReferralLogRequest{
		UserId:       refereeUserID,
		ReferralCode: referralCode,
	}); err != nil {
		return err
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

// calculateKarmaAmount determines the karma amount based on the referral count.
func calculateKarmaAmount(referralCount uint32) uint32 {
	switch {
	case referralCount >= 20:
		return 15000
	case referralCount > 10:
		return 10000
	default:
		return 5000
	}
}
