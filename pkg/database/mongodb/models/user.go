package models

import "github.com/go-playground/validator/v10"

// User struct contains User's table fields
type User struct{
	Name string   `bson:"name" json:"name" validate:"required"`
	Email string   `bson:"email" json:"email" validate:"required"`
	Password string   `bson:"password" json:"password" validate:"required"`
}
// Validate validates the user struct using the validate tag
func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}
