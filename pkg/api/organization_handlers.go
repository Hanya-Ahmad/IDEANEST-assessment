package api

import (
	"net/http"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/api/middleware"
	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/database/mongodb/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrgResponse contains response struct for organization endpoints
type OrgResponse struct{
	Message string `json:"message,omitempty"`
	OrganizationID     string  `json:"organization_id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
	OrganizationMembers []models.Member `json:"organization_members,omitempty"`
}

// UpdateOrgInput contains input struct for /organization/update endpoint
type UpdateOrgInput struct{
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description,omitempty"`
}

// InviteInput contains input struct for /organization/:organization_id/invite endpoint
type InviteInput struct{
	UserEmail string `json:"user_email,omitempty" validate:"required,email"`
}

//  Validate validates the InviteInput struct using the validate tag
func (input *InviteInput) Validate() error {
	validate := validator.New()
	return validate.Struct(input)
}

// createNewOrganization creates a new organization
func (app *App) createNewOrganization(c *gin.Context){
	token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	inputOrg := models.Organization{}
	if err := c.ShouldBindJSON(&inputOrg); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: "invalid json format"})
        return
	}
	if err := inputOrg.Validate(); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: "invalid organization data"})
		return
	}
	id, err := app.DB.CreateNewOrganization(inputOrg)	
	if err == models.ErrOrganizationExists{
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: models.ErrOrganizationExists.Error()})
		return
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message: "failed to create organization"})
        return
 	}
	c.JSON(http.StatusCreated, OrgResponse{ OrganizationID: id})
}

// getOrganizationByID retrieves an organization by its id
func (app *App) getOrganizationByID(c *gin.Context){
	token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	orgID, _ := primitive.ObjectIDFromHex(c.Param("organization_id"))
	org, err := app.DB.GetOrganizationByID(orgID)
	if err ==mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusNotFound, OrgResponse{Message: "organization does not exist"})
        return
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message: "failed to retrieve organization"})
        return
	}
	c.JSON(http.StatusOK, OrgResponse{OrganizationID: org.OrganizationID.Hex(), Name:org.Name, Description: org.Description, OrganizationMembers: org.OrganizationMembers})
}

// getAllOrganizations retrieves all organizations
func (app *App) getAllOrganizations(c *gin.Context){
	token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	organizations, err := app.DB.GetAllOrganizations()
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message:"failed to retrieve organizations"})
		return
	}
		c.JSON(http.StatusOK, organizations)
}

// updateOrganization updates an organization's name and/or description
func (app *App) updateOrganization(c *gin.Context){
	token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	input := UpdateOrgInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: "invalid json format"})
		return
	}
	orgID, _ := primitive.ObjectIDFromHex(c.Param("organization_id"))
	update:= bson.M{"$set": bson.M{"name": input.Name, "description":input.Description}}
	err = app.DB.UpdateOrganization(orgID, update)
	if err ==mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusNotFound, OrgResponse{Message: "organization does not exist"})
        return
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message: "failed to update organization"})
        return
 	}
	org, err := app.DB.GetOrganizationByID(orgID)
	if err ==mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusNotFound, OrgResponse{Message: "organization does not exist"})
        return
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message: "failed to retrieve updated organization"})
        return
 	}
	c.JSON(http.StatusOK, OrgResponse{OrganizationID: org.OrganizationID.Hex(), Name:org.Name, Description: org.Description, OrganizationMembers: org.OrganizationMembers})
}

// deleteOrganization deletes an organization
func (app *App) deleteOrganization(c *gin.Context){
token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	orgID, _ := primitive.ObjectIDFromHex(c.Param("organization_id"))
	err = app.DB.DeleteOrganization(orgID)
	if err ==mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusNotFound, OrgResponse{Message: "organization does not exist"})
        return
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, OrgResponse{Message: "failed to delete organization"})
        return
 	}
	c.JSON(http.StatusOK, OrgResponse{Message: "organization deleted successfully"})

}

// inviteUser invites a user to an organization
func (app *App) inviteUser(c *gin.Context){
	token := c.GetHeader("Authorization")
	_, err := middleware.CheckAuthorization(token, "access")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	input := InviteInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: "invalid json format"})
		return
	}
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, OrgResponse{Message: "invalid user email"})
		return
	}

	orgID, _ := primitive.ObjectIDFromHex(c.Param("organization_id"))
	err =app.DB.SendMail(input.UserEmail,orgID)
	if err ==mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusNotFound, OrgResponse{Message: "organization does not exist"})
        return
	}	
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, OrgResponse{Message: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, OrgResponse{Message: "invitation email sent successfully"})
}

