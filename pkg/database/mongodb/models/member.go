package models

import "github.com/go-playground/validator/v10"

// Member struct contains Member's table fields
type Member struct {
	Name        string `bson:"name" json:"name" validate:"required"`
	Email       string `bson:"email" json:"email" validate:"required"`
	AccessLevel string `bson:"access_level" json:"access_level" validate:"required"`
}

// Validate validates the member struct using the validate tag
func (member *Member) Validate() error {
	validate := validator.New()
	return validate.Struct(member)
}