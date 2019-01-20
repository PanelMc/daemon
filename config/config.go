package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func Init() {
	logrus.WithField("config", "Config").Info("Loading file...")
	viper.SetConfigName("config")
	viper.SetConfigType("hcl")
	viper.AddConfigPath("./")
	if err := viper.ReadInConfig(); err != nil {
		if _, isFile404 := err.(viper.ConfigFileNotFoundError); isFile404 {
			logrus.Warn("Não foi encontrada a config, criando uma nova...")

			if _, err := os.Create("./config.hcl"); err != nil {
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

	//viper.SetDefault(DOCKER_ENDPOINT, "")
	viper.SetDefault(SERVERS_PATH, "./servers")
	viper.SetDefault(DATA_PATH, "./data")

	if err := viper.WriteConfig(); err != nil {
		logrus.Error("Ocurreu um erro ao salvar a config! ", err)
	}

	logrus.Info("Config loaded!")
}
