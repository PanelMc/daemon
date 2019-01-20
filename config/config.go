package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

var config = &Config{}

func GetConfig() *Config {
	return config
}

func Init() {
	logrus.WithField("config", "Config").Info("Loading file...")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, isFile404 := err.(viper.ConfigFileNotFoundError); isFile404 {
			logrus.Warn("Não foi encontrada a config, criando uma nova...")

			if _, err := os.Create("./config.json"); err != nil {
				logrus.WithError(err).Error("Ocurreu um erro ao criar a config! ")
				logrus.Exit(1)
			}
		} else {
			logrus.Error("Ocurreu um erro ao carregar a config! ", err)
			logrus.Exit(1)
		}
	}

	logrus.WithField("config", "Config").Infof("Loaded from file '%s'", viper.ConfigFileUsed())
	logrus.WithField("config", "Config").Info("Loading values...")

	viper.SetDefault("data_path", "./servers")
	viper.SetDefault("file_permissions", 644)
	viper.SetDefault("folder_permissions", 744)
	publicKey, privateKey := generateRsaKeyPair()
	viper.SetDefault("jwt.public_key", string(publicKey))
	viper.SetDefault("jwt.private_key", string(privateKey))

	if err := viper.WriteConfig(); err != nil {
		logrus.Error("Ocurreu um erro ao salvar a config! ", err)
	}

	config.DataPath = viper.GetString("data_path")
	config.FilePermissions = os.FileMode(viper.GetInt32("file_permissions"))
	config.FolderPermissions = os.FileMode(viper.GetInt32("folder_permissions"))

	config.JWT = JwtConfig{
		PublicKey:  []byte(viper.GetString("jwt.public_key")),
		PrivateKey: []byte(viper.GetString("jwt.private_key")),
	}

	logrus.Info("Config loaded!")
}
