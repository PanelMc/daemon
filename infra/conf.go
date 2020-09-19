package infra

import (
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/panelmc/daemon/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitializeLogger initializes the logger with default configurations
func InitializeLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetFormatter(&nested.Formatter{
		HideKeys: true,
		NoColors: true,
	})
	logrus.SetLevel(logrus.TraceLevel)
}

// InitializeCommand initializes the configuration in order to
// load the command flags into it.
func InitializeCommand() {
	logrus.WithField("config", "Config").Info("Loading the configuration file...")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, isFile404 := err.(viper.ConfigFileNotFoundError); isFile404 {
			logrus.Warn("No configuration file found. Creating a new one...")

			if _, err := os.Create("./config.json"); err != nil {
				logrus.WithError(err).Error("Failed to create a new configuration file!")
				logrus.Exit(1)
			}
		} else {
			logrus.Error("Failed to read the configuration file! ", err)
			logrus.Exit(1)
		}
	}

	config.SaveDefaults()

	logrus.WithField("config", "Config").Infof("Loaded the configuration from '%s'", viper.ConfigFileUsed())
}

func ServerConfig(cmd *cobra.Command) {
	cmd.Flags().String("server.host", "127.0.0.1", "host on which the server should listen")
	cmd.Flags().Int("server.port", 8080, "port on which the server should listen")
	cmd.Flags().Bool("server.debug", false, "debug mode for the server")
	cmd.Flags().String("server.allowedOrigins", "*", "allowed origins for the server")

	viper.BindPFlags(cmd.Flags())
}
