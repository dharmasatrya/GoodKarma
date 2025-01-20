package repository

import (
	"context"
	"fmt"
	"os"
	"strings"

	paymentPb "github.com/dharmasatrya/goodkarma/payment-service/proto"
	"github.com/dharmasatrya/goodkarma/user-service/entity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUserSupporter(entity.CreateUserSupporterRequest) (*entity.User, error)
	CreateUserCoordinator(entity.CreateUserCoordinatorRequest) (*entity.User, error)
	Login(entity.LoginRequest) (*entity.User, error)
	GetUserById(string) (*entity.DetailUser, error)
	UpdateProfile(entity.UpdateProfileRequest) (*entity.DetailUser, error)
}

type userRepository struct {
	db *mongo.Database
}

func (ur *userRepository) GetUserCollection() *mongo.Collection {
	return ur.db.Collection("users")
}

func (ur *userRepository) GetProfileCollection() *mongo.Collection {
	return ur.db.Collection("profiles")
}

func NewUserRepository(DB *mongo.Database) UserRepository {
	return &userRepository{
		db: DB,
	}
}

func (ur *userRepository) CreateUserSupporter(request entity.CreateUserSupporterRequest) (*entity.User, error) {
	userCollection := ur.GetUserCollection()
	profileCollection := ur.GetProfileCollection()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	if err := ur.validateCreateUser(request); err != nil {
		return nil, err
	}

	newUser := entity.User{
		ID:       primitive.NewObjectID(),
		Username: request.Username,
		Email:    request.Email,
		Password: string(hashedPassword),
		Role:     request.Role,
	}

	insertUser, err := userCollection.InsertOne(context.Background(), newUser)

	if err != nil {
		return nil, err
	}

	profile := entity.Profile{
		ID:       primitive.NewObjectID(),
		UserID:   insertUser.InsertedID.(primitive.ObjectID),
		FullName: request.FullName,
		Address:  request.Address,
		Phone:    request.Phone,
		Photo:    request.Photo,
	}

	_, err = profileCollection.InsertOne(context.Background(), profile)

	if err != nil {
		return nil, err
	}

	return &newUser, nil
}

func (ur *userRepository) CreateUserCoordinator(request entity.CreateUserCoordinatorRequest) (*entity.User, error) {
	var user *entity.CreateUserSupporterRequest

	user = &entity.CreateUserSupporterRequest{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
		Role:     request.Role,
		FullName: request.FullName,
		Address:  request.Address,
		Phone:    request.Phone,
		Photo:    request.Photo,
	}

	res, err := ur.CreateUserSupporter(*user)

	if err != nil {
		return nil, err
	}

	err = createWallet(res.ID.Hex(), request)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (ur *userRepository) Login(request entity.LoginRequest) (*entity.User, error) {
	userCollection := ur.GetUserCollection()

	var user entity.User
	var filter primitive.M

	if strings.Contains(request.UsernameOrEmail, "@") {
		filter = primitive.M{"email": request.UsernameOrEmail}
	} else {
		filter = primitive.M{"username": request.UsernameOrEmail}
	}

	err := userCollection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserById(id string) (*entity.DetailUser, error) {
	userCollection := ur.GetUserCollection()
	profileCollection := ur.GetProfileCollection()

	var user entity.DetailUser
	var profile entity.Profile

	userID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	err = userCollection.FindOne(context.Background(), primitive.M{"_id": userID}).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	err = profileCollection.FindOne(context.Background(), primitive.M{"user_id": userID}).Decode(&profile)

	if err != nil {
		return nil, err
	}

	user.FullName = profile.FullName
	user.Address = profile.Address
	user.Phone = profile.Phone
	user.Photo = profile.Photo

	return &user, nil
}

func (ur *userRepository) UpdateProfile(request entity.UpdateProfileRequest) (*entity.DetailUser, error) {
	profileCollection := ur.GetProfileCollection()

	userID, err := primitive.ObjectIDFromHex(request.UserID)

	if err != nil {
		return nil, err
	}

	_, err = ur.GetUserById(request.UserID)

	if err != nil {
		return nil, err
	}

	update := primitive.M{
		"$set": primitive.M{
			"full_name": request.FullName,
			"address":   request.Address,
			"phone":     request.Phone,
			"photo":     request.Photo,
		},
	}

	_, err = profileCollection.UpdateOne(context.Background(), primitive.M{"user_id": userID}, update)

	if err != nil {
		return nil, err
	}

	user, err := ur.GetUserById(request.UserID)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func createWallet(userID string, request entity.CreateUserCoordinatorRequest) error {
	// Get payment service URI
	paymentServiceURI := os.Getenv("PAYMENT_SERVICE_URI")
	if paymentServiceURI == "" {
		return fmt.Errorf("PAYMENT_SERVICE_URI is not set")
	}

	// Create gRPC connection
	grpcConn, err := grpc.NewClient(paymentServiceURI, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to payment service: %w", err)
	}
	defer grpcConn.Close()

	// Initialize payment client
	paymentClient := paymentPb.NewPaymentServiceClient(grpcConn)

	// Make the gRPC call to create a wallet
	_, err = paymentClient.CreateWallet(context.Background(), &paymentPb.CreateWalletRequest{
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

func (ur *userRepository) validateCreateUser(request entity.CreateUserSupporterRequest) error {
	userCollection := ur.GetUserCollection()

	if checkUsername, _ := userCollection.CountDocuments(context.Background(), primitive.M{"username": request.Username}); checkUsername > 0 {
		return fmt.Errorf("username %v already exists", request.Username)
	}

	if checkEmail, _ := userCollection.CountDocuments(context.Background(), primitive.M{"email": request.Email}); checkEmail > 0 {
		return fmt.Errorf("email %v already exists", request.Email)
	}

	if checkPhone, _ := userCollection.CountDocuments(context.Background(), primitive.M{"phone": request.Phone}); checkPhone > 0 {
		return fmt.Errorf("phone %v already exists", request.Phone)
	}

	return nil
}
