package daemon

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/panelmc/daemon/types"

	"github.com/panelmc/daemon/config"
	"github.com/sirupsen/logrus"
)

func Init() {
	serversPath := DataPath()
	if err := os.MkdirAll(serversPath, config.GetConfig().FolderPermissions); err != nil {
		logrus.WithError(err).Fatal("Failed to create de servers data directory!")
	}
	logrus.Infof("Loading servers from '%s'...", serversPath)

	if files, err := ioutil.ReadDir(serversPath); err == nil {
		for _, f := range files {
			// Get all servers folders
			if f.IsDir() {
				logrus.Debugf("Checking the directory '%s'", f.Name())

				// Check it theres a config file present
				cFile := filepath.Join(serversPath, f.Name(), "config", "config.json")

				if _, err := os.Stat(cFile); err == nil {
					logrus.WithField("server", f.Name()).Debug("Found a config file, checking...")

					loadConfig(f.Name(), cFile)
				} else if os.IsNotExist(err) {
					logrus.Debugf("The file '%s' wasn't found, skipping.", cFile)
				} else {
					logrus.WithError(err).Error("Failed to check if config file is present.")
				}
			}
		}
	} else {
		log.Fatal(err)
	}

	logrus.Infof("Loaded a total of %d servers!", len(servers))
}

func loadConfig(rawName, configPath string) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.WithField("server", rawName).WithError(err).
			Error("There was an error while trying to read the config file!")
		return
	}

	serverConfig := &types.ServerConfiguration{}
	if err := json.Unmarshal([]byte(file), serverConfig); err != nil {
		logrus.WithField("server", rawName).WithError(err).
			Debugf("The 'config.json' is not a valid server configuration!")
		return
	}

	logrus.WithField("server", serverConfig.ID).
		Debug("Config loaded, initializing the server...")

	server := &Server{
		ID:       serverConfig.ID,
		Name:     serverConfig.Name,
		Type:     serverConfig.Type,
		Settings: serverConfig.Settings,
		Container: DockerContainer{
			ContainerID: serverConfig.Container.ContainerID,
			Image:       serverConfig.Container.Image,
		},
	}

	if err := server.Init(); err != nil {
		logrus.WithField("server", server.ID).WithError(err).
			Errorf("There was an error while initializing the server %s", server.ID)
	} else {
		servers[server.ID] = server
		logrus.WithField("server", server.ID).
			Infof("The server '%s' was successfuly initialized!", server.Name)
	}
}

func Start() {
	for _, server := range servers {
		if err := server.Start(); err != nil {
			logrus.WithField("server", server.ID).WithError(err).Error("There was an error while trying to start the server.")
		}
	}
}
