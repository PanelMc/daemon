package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/panelmc/daemon/config"
	"github.com/sirupsen/logrus"
)

type TokenClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// NewToken generates a mock token
func NewToken() string {
	claims := TokenClaims{
		UserID:   "heroslender",
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
