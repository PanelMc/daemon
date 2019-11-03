package daemon

import (
	"github.com/panelmc/daemon/types"
	"github.com/sirupsen/logrus"
)

type IServer interface {
	Start() error
	Stop() error
	Restart() error

	Execute(command string) error

	Init() error
	DataPath() string
	Save() error
}

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	willRestart bool

	Settings types.ServerSettings `json:"settings"`

	Stats *types.ServerStats `json:"stats,omitempty"`

	Container DockerContainer `json:"container"`
}

var _ IServer = &Server{}

func (s *Server) Init() error {
	if err := NewDockerContainer(s); err != nil {
		return err
	}

	if s.Container.ContainerID == "" {
		logrus.WithField("server", s.ID).Debug("Creating a new container...")

		if err := s.Container.Create(); err != nil {
			logrus.
				WithField("server", s.ID).
				WithError(err).
				Error("Failed to create the docker container.")
			return err
		}
	}

	return nil
}
