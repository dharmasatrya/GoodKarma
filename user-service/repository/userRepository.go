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
	CreateUserSupporter(entity.CreateUserSupporterRequest) (*entity.User, error)
	CreateUserCoordinator(entity.CreateUserCoordinatorRequest) (*entity.User, error)
	Login(entity.LoginRequest) (*entity.User, error)
	GetUserById(string) (*entity.DetailUser, error)
	UpdateProfile(entity.UpdateProfileRequest) (*entity.DetailUser, error)
}

type userRepository struct {
	db             *mongo.Database
	userRepository UserRepository
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
		ID:            primitive.NewObjectID(),
		Username:      request.Username,
		Email:         request.Email,
		Password:      string(hashedPassword),
		Role:          request.Role,
		EmailVerified: false,
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

	// err = createWallet(res.ID.Hex(), request)

	// if err != nil {
	// 	return nil, err
	// }

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

	// Check if email is verified
	if user.EmailVerified == false {
		return nil, fmt.Errorf("email not verified")
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
