package daemon

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/panelmc/daemon/config"
	"github.com/panelmc/daemon/types"
	"github.com/sirupsen/logrus"
)

func (server *Server) Save() error {
	if err := os.MkdirAll(server.DataPath(), config.GetConfig().FolderPermissions); err != nil {
		return err
	}

	// Prevent stats from saving to config, but keep them on the API
	toSave := &types.ServerConfiguration{
		ID:       server.ID,
		Name:     server.Name,
		Type:     server.Type,
		Settings: server.Settings,
		Container: types.DockerContainerConfiguration{
			ContainerID: server.Container.ContainerID,
			Image:       server.Container.Image,
		},
	}

	serverJSON, err := json.MarshalIndent(toSave, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(server.ConfigFilePath(), serverJSON, config.GetConfig().FilePermissions); err != nil {
		return err
	}
	return nil
}

func (s *Server) DataPath() string {
	return filepath.Join(DataPath(), s.ID)
}

func (s *Server) ConfigFilePath() string {
	return filepath.Join(s.DataPath(), "config", "config.json")
}

func DataPath() string {
	dataPath, err := filepath.Abs(config.GetConfig().DataPath)
	if err != nil {
		logrus.WithError(err).Error("Failed to get absolute data path for server.")
		dataPath = config.GetConfig().DataPath
	}

	return filepath.Join(dataPath)
}
