package middleware

import (
	"strings"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/utils"
	"github.com/rs/zerolog/log"
)

// CheckAuthorization checks the authorization token and returns token claims
func CheckAuthorization(token string, tokenType string) (*utils.TokenClaims, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil, utils.ErrMissingToken
	}
	tokenClaims, err := utils.ValidateToken(token, tokenType)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, utils.ErrInvalidToken
	}
	return tokenClaims, nil
}