package api

import (
	"net/http"
	"time"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/api/middleware"
	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/database/mongodb/models"
	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserResponse contains the response returned in user endpoints
type UserResponse struct {
    Message string `json:"message"`
	AccessToken string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// createNewUser signs up a new user
func (app *App) createNewUser(c *gin.Context){
	user := models.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid json format"})
        return
	}
	if err := user.Validate(); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid user data"})
		return
	}
	hashedPassword, err := utils.HashPassword([]byte(user.Password)) 
	if err!= nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, UserResponse{Message:"failed to encrypt password"})
		return
	}
	user.Password = string(hashedPassword)
 	err = app.DB.CreateNewUser(user)
	if err == models.ErrUserExists{
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: models.ErrUserExists.Error()})
		return
	}
 	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to create user"})
        return
 	}
 	c.JSON(http.StatusCreated, UserResponse{Message: "user created successfully"})
}

// signIn signs in an existing user
func (app *App) signIn(c *gin.Context){
	input := models.SignInInput{}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid json format"})
		return
	}	
	if err := input.Validate(); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid user data"})
		return
	}
	user, err := app.DB.GetUserByEmail(input.Email)
	if err == mongo.ErrNoDocuments{
		log.Error().Err(err).Send()
		c.JSON(http.StatusCreated, UserResponse{Message: models.ErrUserNotFound.Error()})
		return			
	}
	if err!=nil{
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to retrieve user from database"})
        return
 	}
	ok := utils.CheckPasswordHash([]byte(user.Password),input.Password)
	if !ok{
		log.Error().Err(err).Send()
 		c.JSON(http.StatusCreated, UserResponse{Message: "incorrect password"})
		return
	}
	accessToken, err := utils.GenerateToken(user.ID, time.Now().Add(time.Minute*15),"access")
	if err!=nil{
		log.Error().Err(err).Send()
 		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to create access token"})
		return
	}
	refreshToken, err := utils.GenerateToken(user.ID, time.Now().Add(time.Hour*24),"refresh")
	if err!=nil{
		log.Error().Err(err).Send()
 		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to create refresh token"})
		return
	}
 	c.JSON(http.StatusCreated, UserResponse{Message: "signin successful", AccessToken: accessToken, RefreshToken: refreshToken})

}

// refreshToken creates new access and refresh tokens for a signed-in user
func (app *App) refreshToken(c *gin.Context){
	var input models.RefreshToken
	
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid json format"})
		return
	}		

	if err := input.Validate(); err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, UserResponse{Message: "invalid token data"})
		return
	}
	token := input.Token
	tokenClaims, err := middleware.CheckAuthorization(token, "refresh")
	if err != nil {
		log.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, UserResponse{Message: err.Error()})
		return
	}

	accessToken, err := utils.GenerateToken(tokenClaims.ID, time.Now().Add(time.Minute*15),"access")
	if err!=nil{
		log.Error().Err(err).Send()
 		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to create access token"})
		return
	}
	refreshToken, err := utils.GenerateToken(tokenClaims.ID, time.Now().Add(time.Hour*24),"refresh")
	if err!=nil{
		log.Error().Err(err).Send()
 		c.JSON(http.StatusInternalServerError, UserResponse{Message: "failed to create refresh token"})
		return
	}
 	c.JSON(http.StatusCreated, UserResponse{Message: "tokens refreshed successfully", AccessToken: accessToken, RefreshToken: refreshToken})
}