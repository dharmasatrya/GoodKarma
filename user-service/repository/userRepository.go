package repository

import (
	"context"
	"goodkarma-user-service/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	CreateUser(request entity.CreateUserRequest) (*entity.User, error)
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

func (ur *userRepository) CreateUser(request entity.CreateUserRequest) (*entity.User, error) {
	userCollection := ur.GetUserCollection()
	profileCollection := ur.GetProfileCollection()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)

	if err != nil {
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
