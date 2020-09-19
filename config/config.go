package config

import (
	"os"

	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var config *Config

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}

	return config
}

func SaveDefaults() {
	viper.SetDefault("data_path", "./servers")
	viper.SetDefault("file_permissions", 644)
	viper.SetDefault("folder_permissions", 744)

	if !viper.IsSet("jwt.public_key") || !viper.IsSet("jwt.private_key") {
		publicKey, privateKey := generateRsaKeyPair()
		viper.Set("jwt.public_key", string(publicKey))
		viper.Set("jwt.private_key", string(privateKey))
	}

	if err := viper.WriteConfig(); err != nil {
		logrus.Error("Failed to save the configuration! ", err)
	}
}

func NewConfig() *Config {
	return &Config{
		DataPath:          viper.GetString("data_path"),
		FilePermissions:   os.FileMode(viper.GetInt32("file_permissions")),
		FolderPermissions: os.FileMode(viper.GetInt32("folder_permissions")),
		JWT: JwtConfig{
			PublicKey:  []byte(viper.GetString("jwt.public_key")),
			PrivateKey: []byte(viper.GetString("jwt.private_key")),
		},
	}
}

func (c *Config) Save() {
	viper.Set("data_path", c.DataPath)
	viper.Set("file_permissions", c.FilePermissions)
	viper.Set("folder_permissions", c.FolderPermissions)
	viper.Set("jwt.public_key", string(c.JWT.PublicKey))
	viper.Set("jwt.private_key", string(c.JWT.PrivateKey))

	if err := viper.WriteConfig(); err != nil {
		logrus.Error("Failed to save the configuration! ", err)
	}
}

// generateRsaKeyPair generates RSA public and private keys, PEM encoded
func generateRsaKeyPair() ([]byte, []byte) {
	reader := rand.Reader
	bitSize := 2048
	key, err := rsa.GenerateKey(reader, bitSize)
	if err != nil {
		logrus.WithField("auth", "jwt").WithError(err).Error("Failed to generate the RSA Keys")
		return nil, nil
	}
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		logrus.WithField("auth", "jwt").WithError(err).Error("Failed to marshal public key")
		return nil, nil
	}
	var publickey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var privateKey = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	public := pem.EncodeToMemory(publickey)
	private := pem.EncodeToMemory(privateKey)
	if public == nil || private == nil {
		logrus.WithField("auth", "jwt").Error("Failed to encode the RSA Keys")
		return nil, nil
	}
	return public, private
}
