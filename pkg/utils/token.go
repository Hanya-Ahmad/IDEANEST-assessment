package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ErrInvalidToken indicates that a token is invalid
var ErrInvalidToken = errors.New("failed to validate token")

// ErrMissingToken indicates that no bearer token was included
var ErrMissingToken = errors.New("token is required")

// TokenClaims contains JWT token claims
type TokenClaims struct {
	ID primitive.ObjectID
	jwt.RegisteredClaims
}

// GenerateToken creates new access and refresh tokens
func GenerateToken(userID primitive.ObjectID, expirationTime time.Time, tokenType string) (string, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	tokenClaims := TokenClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	method := jwt.SigningMethodHS256
	if tokenType == "refresh" {
		method = jwt.SigningMethodHS384 // Use a different signing method for refresh tokens, if desired
	}

	token := jwt.NewWithClaims(method, tokenClaims)
	return token.SignedString([]byte(secret))
}

// ValidateToken checks the validity of a token
func ValidateToken(signedToken string, tokenType string) (*TokenClaims, error) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	secret := os.Getenv("JWT_SECRET")
	method := jwt.GetSigningMethod("HS256")
	if tokenType == "refresh" {
		method = jwt.GetSigningMethod("HS384") // Use a different signing method for refresh tokens, if you've changed it in GenerateToken
	}
	token, err := jwt.ParseWithClaims(signedToken, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if token.Method != method {
		return nil,ErrInvalidToken
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	if time.Now().After((*claims.ExpiresAt).Time) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
