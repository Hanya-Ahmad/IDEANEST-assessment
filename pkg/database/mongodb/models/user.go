package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrUserExists is raised when a new user is added with an email linked to a registered user
var ErrUserExists = errors.New("a user with that email already exists")

// ErrUserNotFound is raised when a nonexisting user tries to sign in
var ErrUserNotFound = errors.New("this email is not registered to an account")


// User struct contains User's table fields
type User struct{
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string   `bson:"name" json:"name,omitempty" validate:"required"`
	Email string   `bson:"email" json:"email" validate:"required,email"`
	Password string   `bson:"password" json:"password" validate:"required"`
}

// SignInInput contains user's signin input
type SignInInput struct{
	Email string   `bson:"email" json:"email" validate:"required,email"`
	Password string   `bson:"password" json:"password" validate:"required"`	
}

// RefreshToken contains the /refresh-token request body
type RefreshToken struct{
	Token string  `bson:"refresh_token" json:"refresh_token" validate:"required"`
}

// Validate validates the user struct using the validate tag
func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}
 
// Validate validates the user struct using the validate tag
func (user *SignInInput) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}

// Validate validates the RefreshToken struct using the validate tag
func (token *RefreshToken) Validate() error {
	validate := validator.New()
	return validate.Struct(token)
}

// CreateNewUser adds a new user to the database
func (db *DBClient) CreateNewUser(user User) error{
	collection := db.Client.Collection("User")
	_, err := db.GetUserByEmail(user.Email)
	if err == mongo.ErrNoDocuments{
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	newID := result.InsertedID
	fmt.Println("new id: ", newID)
	return nil
	}else{
		return ErrUserExists
	}
}

// GetUserByEmail retrieves a user from the database by their email
func (db *DBClient) GetUserByEmail(email string)(User, error){
	collection := db.Client.Collection("User")
	filter := bson.M{"email": email}
	var checkedUser User
    err := collection.FindOne(context.Background(), filter).Decode(&checkedUser)
	return checkedUser, err
}
