package daemon

import (
	"encoding/json"
	"github.com/heroslender/panelmc/config"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type ServerMap map[string]*ServerStruct

var servers = make(ServerMap)

func GetServers() *ServerMap {
	return &servers
}

// Get a server by it's id. If no server found, return nil
func GetServerById(id string) *ServerStruct {
	for _, s := range servers {
		if s.Id == id {
			return s
		}
	}

	return nil
}

// Get a server by it's id or name. If no server found, return nil
func GetServer(server string) *ServerStruct {
	for _, s := range servers {
		if s.Id == server || s.Name == server {
			return s
		}
	}

	return nil
}

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
					file, err := ioutil.ReadFile(cFile)
					if err != nil {
						logrus.WithField("server", f.Name()).WithError(err).
							Error("There was an error while trying to read the config file!")
						continue
					}

					server := &ServerStruct{}
					if err := json.Unmarshal([]byte(file), server); err != nil {
						logrus.WithField("server", f.Name()).WithError(err).
							Debugf("The 'config.json' is not a valid server configuration!")
						continue
					}

					logrus.WithField("server", server.Id).
						Debug("Config loaded, initializing the server...")
					if err := server.Init(); err != nil {
						logrus.WithField("server", server.Id).WithError(err).
							Errorf("There was an error while initializing the server %s", server.Id)
					} else {
						servers[server.Id] = server
						logrus.WithField("server", server.Id).
							Infof("The server '%s' was successfuly initialized!", server.Name)
					}
				} else if !os.IsNotExist(err) {
					logrus.WithError(err).Error("Failed to check if config file is present.")
				} else {
					logrus.Debugf("The file '%s' wasn't found, skipping.", cFile)
				}
			}
		}
	} else {
		log.Fatal(err)
	}

	logrus.Infof("Loaded a total of %d servers!", len(servers))
}

func Start() {
	for _, server := range servers {
		if err := server.Start(); err != nil {
			logrus.WithField("server", server.Id).WithError(err).Error("There was an error while trying to start the server.")
		}
	}
}

func saveServerConfig(server *ServerStruct) error {
	if err := os.MkdirAll(server.DataPath(), config.GetConfig().FolderPermissions); err != nil {
		return err
	}


	// Prevent stats from saving to config, but keep them on the API
	toSave := *server
	toSave.Stats = nil

	serverJSON, err := json.MarshalIndent(&toSave, "", "  ")
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(server.ConfigFilePath(), serverJSON, config.GetConfig().FilePermissions); err != nil {
		return err
	}
	return nil
}

func (s *ServerStruct) Save() error {
	return saveServerConfig(s)
}
