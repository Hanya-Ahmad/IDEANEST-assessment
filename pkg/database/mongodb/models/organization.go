package models

import "github.com/go-playground/validator/v10"

// Organization struct contains Organization's table fields
type Organization struct{
	// OrganizationID     string   `bson:"organization_id,omitempty" json:"organization_id"`
	Name               string   `bson:"name" json:"name"`
	Description        string   `bson:"description" json:"description"`
	OrganizationMembers []Member `bson:"organization_members" json:"organization_members"`
}
// Validate validates the organization struct using the validate tag
func (organization *Organization) Validate() error {
	validate := validator.New()
	return validate.Struct(organization)
}