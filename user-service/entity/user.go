package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Username string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Role     string             `json:"role" bson:"role"`
}

type Profile struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	UserID   primitive.ObjectID `json:"user_id" bson:"user_id"`
	FullName string             `json:"full_name" bson:"full_name"`
	Address  string             `json:"address" bson:"address"`
	Phone    string             `json:"phone" bson:"phone"`
	Photo    string             `json:"photo" bson:"photo"`
}

type Wallet struct {
	ID                primitive.ObjectID `json:"id" bson:"_id"`
	UserID            primitive.ObjectID `json:"user_id" bson:"user_id"`
	BankAccountName   string             `json:"bank_account_name" bson:"bank_account_name"`
	BankCode          string             `json:"bank_code" bson:"bank_code"`
	BankAccountNumber string             `json:"bank_account_number" bson:"bank_account_number"`
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
