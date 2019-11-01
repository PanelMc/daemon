package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/config"
	"github.com/sirupsen/logrus"
)

// GinHandler - JWT check for gin middleware
func GinHandler(c *gin.Context) {
	if err := VerifyRequest(c.Request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"error": "unauthorized",
				//"message": "You need to login to access this content.",
				// TODO remove for produtction
				"message": err.Error(),
			},
		})
	}
}

// SocketHandler - JWT check for the socket middleware
func SocketHandler(r *http.Request) (http.Header, error) {
	return nil, VerifyRequest(r)
}

// VerifyRequest - Verify the request for valid token
func VerifyRequest(r *http.Request) error {
	token := r.Header.Get("Authorization")
	if token == "" {
		// If not found on headers, get from url query
		token = r.URL.Query().Get("Authorization")
	}

	if token == "" {
		return errors.New("Required authorization token not found")
	}
	token = token[7:]

	if token, err := Verify(token); err != nil {
		logrus.WithField("auth", "jwt").WithError(err).Error("Failed to verify Token!")
		return err
	} else {
		newRequest := r.WithContext(context.WithValue(r.Context(), "jwt", token.Claims.(jwt.MapClaims)))
		// UpdateStats the current request with the new context information.
		*r = *newRequest
	}

	return nil
}

// Verify and parse the token string
func Verify(token string) (*jwt.Token, error) {
	if token == "" {
		return &jwt.Token{}, errors.New("Required authorization token not found")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(config.GetConfig().JWT.PublicKey)
	})

	if err != nil {
		return parsedToken, fmt.Errorf("Error parsing token: %v", err)
	}

	if jwt.SigningMethodRS256.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			jwt.SigningMethodRS256.Alg(),
			parsedToken.Header["alg"])
		return parsedToken, fmt.Errorf("Error validating token algorithm: %s", message)
	}

	if !parsedToken.Valid {
		return parsedToken, errors.New("Token is invalid")
	}

	return parsedToken, nil
}
