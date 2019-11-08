package daemon

import (
	"github.com/panelmc/daemon/types"
)

func (s *Server) Start() error {
	if err := s.Container.Start(); err != nil {
		return err
	}

	s.UpdateStatus(types.ServerStatusStarting)
	return nil
}

func (s *Server) Stop() error {
	if err := s.Container.Stop(); err != nil {
		return err
	}

	s.UpdateStatus(types.ServerStatusStopping)
	return nil
}

func (s *Server) Restart() error {
	if err := s.Stop(); err != nil {
		return err
	}
	if err := s.Start(); err != nil {
		return err
	}

	return nil
}

func (s *Server) Execute(command string) error {
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

	s.Logger().Infof("Executed console command /%s", command)
	return nil
}

func (s *Server) UpdateStatus(status types.ServerStatus) {
	if s.Stats.Status != status {
		s.Stats.Status = status
		s.onStatusUpdate(status)
	}
}
