package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/dharmasatrya/goodkarma/user-service/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(entity.CreateUserRequest) (*entity.User, error)
	CreateMerchant(entity.CreateUserRequest, entity.CreateMerchantRequest) (*entity.User, error)
	Login(entity.LoginRequest) (*entity.User, error)
	GetUserById(string) (*entity.DetailUser, error)
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

func (ur *userRepository) GetWalletCollection() *mongo.Collection {
	return ur.db.Collection("wallets")
}

func NewUserRepository(DB *mongo.Database) UserRepository {
	return &userRepository{
		db: DB,
	}
}

func (ur *userRepository) CreateUser(request entity.CreateUserRequest) (*entity.User, error) {
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

func (ur *userRepository) CreateMerchant(request entity.CreateUserRequest, requestMerchant entity.CreateMerchantRequest) (*entity.User, error) {
	walletCollection := ur.GetWalletCollection()

	user, err := ur.CreateUser(request)

	if err != nil {
		return nil, err
	}

	wallet := entity.Wallet{
		ID:                primitive.NewObjectID(),
		UserID:            user.ID,
		AccountHolderName: requestMerchant.AccountHolderName,
		BankCode:          requestMerchant.BankCode,
		BankAccountNumber: requestMerchant.BankAccountNumber,
	}

	_, err = walletCollection.InsertOne(context.Background(), wallet)

	if err != nil {
		return nil, err
	}

	return user, nil
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

	fmt.Println(user)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) GetUserById(id string) (*entity.DetailUser, error) {
	userCollection := ur.GetUserCollection()
	profileCollection := ur.GetProfileCollection()

	var user entity.DetailUser

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

	var profile entity.Profile

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

func (ur *userRepository) validateCreateUser(request entity.CreateUserRequest) error {
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
