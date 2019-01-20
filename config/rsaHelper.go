package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/sirupsen/logrus"
)

// Generate RSA public and private keys, PEM encoded
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