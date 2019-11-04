package daemon

import (
	"context"

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

	if err := s.Container.prepare(context.Background()); err != nil {
		logrus.
			WithField("server", s.ID).
			WithError(err).
			Error("Failed to prepare the container for the server.")
		return err
	}

	return nil
}
