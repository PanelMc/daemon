package config

import "os"

type Config struct {
	// Folder to store all server files
	DataPath          string      `mapstructure:"data_path"`
	// Permissions to set when creating new files
	FilePermissions   os.FileMode `mapstructure:"file_permissions"`
	// Permissions to set when creating new folders
	FolderPermissions os.FileMode `mapstructure:"folder_permissions"`

	// JWT Properties
	JWT JwtConfig `json:"jwt"`
}

type JwtConfig struct {
	// RSA Public key PEM encoded
	PublicKey  []byte `mapstructure:"public_key"`
	// RSA Private key PEM encoded
	PrivateKey []byte `mapstructure:"private_key"`
}