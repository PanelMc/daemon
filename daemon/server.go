package daemon

import (
	"github.com/heroslender/panelmc/api/socket"
	"github.com/heroslender/panelmc/config"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"time"
)

func (s *ServerStruct) Start() error {
	if err := s.Container.Start(); err != nil {
		return err
	}

	s.UpdateStatus(ServerStatusStarting)
	return nil
}

func (s *ServerStruct) Stop() error {
	if err := s.Container.Stop(); err != nil {
		return err
	}

	s.UpdateStatus(ServerStatusStopping)
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
	dataPath, err := filepath.Abs(config.GetConfig().DataPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get absolute data path for server.")
		dataPath = config.GetConfig().DataPath
	}

	return filepath.Join(dataPath)
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

func (s *ServerStruct) onDie() {
	s.UpdateStatus(ServerStatusOffline)

	if s.willRestart {
		payload := socket.ServerConsolePayload{
			ServerId: s.Id,
			Line:     "Server stopped! Restarting in 5 seconds...",
		}
		logrus.WithField("server", s.Id).WithField("event", "Console").Info(payload.Line)
		socket.BroadcastTo(s.Id, "console_output", payload)

		go func() {
			time.Sleep(5 * time.Second)
			if s.Stats.Status == ServerStatusOffline {
				if err := s.Start(); err != nil {
					logrus.WithField("server", s.Id).WithError(err).Error("There was an error while trying to start the server.")
				}
			}
		}()
	}
}
