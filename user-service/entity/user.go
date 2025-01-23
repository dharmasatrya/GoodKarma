package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID            primitive.ObjectID `json:"id" bson:"_id"`
	Username      string             `json:"username" bson:"username"`
	Email         string             `json:"email" bson:"email"`
	Password      string             `json:"password" bson:"password"`
	Role          string             `json:"role" bson:"role"`
	EmailVerified bool               `json:"email_verified" bson:"email_verified"`
	ReferralCode  string             `json:"referral_code" bson:"referral_code"`
}

type Profile struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	FullName string             `json:"full_name" bson:"full_name"`
	Address  string             `json:"address" bson:"address"`
	Phone    string             `json:"phone" bson:"phone"`
	Photo    string             `json:"photo" bson:"photo"`
}

type DetailUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Role     string             `json:"role" bson:"role"`
	FullName string             `json:"full_name" bson:"full_name"`
	Address  string             `json:"address" bson:"address"`
	Phone    string             `json:"phone" bson:"phone"`
	Photo    string             `json:"photo" bson:"photo"`
}

type CreateUserSupporterRequest struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Role         string `json:"role"`
	FullName     string `json:"full_name"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Photo        string `json:"photo"`
	ReferralCode string `json:"referral_code"`
}

type CreateUserCoordinatorRequest struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	Role              string `json:"role"`
	FullName          string `json:"full_name"`
	Address           string `json:"address"`
	Phone             string `json:"phone"`
	Photo             string `json:"photo"`
	AccountHolderName string `json:"account_holder_name"`
	BankCode          string `json:"bank_code"`
	BankAccountNumber string `json:"bank_account_number"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}

type CreateMerchantRequest struct {
	AccountHolderName string `json:"account_holder_name"`
	BankCode          string `json:"bank_code"`
	BankAccountNumber string `json:"bank_account_number"`
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UpdateProfileRequest struct {
	UserID   string `json:"user_id"`
	FullName string `json:"full_name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Photo    string `json:"photo"`
}
