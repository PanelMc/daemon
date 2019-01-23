package jwt

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/heroslender/panelmc/config"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type TokenClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewToken() string {
	claims := TokenClaims{
		UserId:   "heroslender",
		Username: "Heroslender",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	key, err := jwt.ParseRSAPrivateKeyFromPEM(config.GetConfig().JWT.PrivateKey)
	if err != nil {
		logrus.WithField("auth", "jwt").WithError(err).Error("Failed to parse the RSA Private Key!")
		return ""
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		logrus.WithField("auth", "jwt").WithError(err).Error("Failed to generate a new Token!")
		return ""
	}

	return tokenString
}

func GinHandler(c *gin.Context) {
	if err := VerifyRequest(c.Request); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": gin.H{
				"error":   "unauthorized",
				//"message": "You need to login to access this content.",
				// TODO remove for produtction
				"message": err.Error(),
			},
		})
	}
}

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
		newRequest := r.WithContext(context.WithValue(r.Context(), "jwt", token.Claims.(*TokenClaims)))
		// UpdateStats the current request with the new context information.
		*r = *newRequest
	}

	return nil
}

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
