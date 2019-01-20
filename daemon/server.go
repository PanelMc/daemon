package daemon

import (
	"github.com/heroslender/panelmc/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
)

func (s *ServerStruct) Start() error {
	if err := s.Container.Start(); err != nil {
		return err
	}
	return nil
}

func (s *ServerStruct) Stop() error {
	if err := s.Container.Stop(); err != nil {
		return err
	}

	return nil
}

func (s *ServerStruct) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	if err := s.Start(); err != nil {
		return err
	}

	return nil
}

func (s *ServerStruct) Execute(command string) error {
	switch command {
	case "stop":
		if err := s.Stop(); err != nil {
			return err
		}
		break
	case "restart":
		if err := s.Start(); err != nil {
			return err
		}
		break
	default:
		if err := s.Container.Exec(command); err != nil {
			return err
		}
	}

	logrus.WithField("server", s.Name).Infof("Executed console command /%s", command)
	return nil
}

func DataPath() string {
	dataPath, err := filepath.Abs(viper.GetString(config.DATA_PATH))
	if err != nil {
		logrus.WithError(err).Error("Failed to get absolute data path for server.")
		dataPath = viper.GetString(config.DATA_PATH)
	}

	return filepath.Join(dataPath, viper.GetString(config.SERVERS_PATH))
}

func (s *ServerStruct) DataPath() string {
	return filepath.Join(DataPath(), s.Id)
}

func (s *ServerStruct) ConfigFilePath() string {
	return filepath.Join(s.DataPath(), "config", "config.json")
}

func (s *ServerStruct) Init() error {
	if err := NewDockerContainer(s); err != nil {
		return err
	}

	if s.Container.ContainerId == "" {
		logrus.WithField("server", s.Id).Debug("Creating a new container...")
		if err := s.Container.Create(); err != nil {
			logrus.WithField("server", s.Id).Error("Failed to create the docker container.")
			return err
		}
	}
	return nil
}
