package daemon

import (
	"strings"

	"github.com/panelmc/daemon/types"
)

// ServerMap is an alias for map[string]*Server
type ServerMap map[string]*Server

var servers = make(ServerMap)

// GetServers - Get all servers avaliable
func GetServers() *ServerMap {
	return &servers
}

// GetServerByID - Get a server by it's id. If no server found, return nil
func GetServerByID(id string) *Server {
	for _, s := range servers {
		if s.ID == id {
			return s
		}
	}

	return nil
}

// GetServer - Get a server by it's id or name. If no server found, return nil
func GetServer(server string) *Server {
	for _, s := range servers {
		if s.ID == server || strings.EqualFold(s.Name, server) {
			return s
		}
	}

	return nil
}

// NewServer - Create a new Server object based on the given configuration
func NewServer(config *types.ServerConfiguration) *Server {
	if config.ID == "" {
		config.ID = strings.ReplaceAll(config.Name, " ", "_")
	}

	server := &Server{
		ID:       config.ID,
		Name:     config.Name,
		Type:     config.Type,
		Settings: config.Settings,
		Container: DockerContainer{
			ContainerID: config.Container.ContainerID,
			Image:       config.Container.Image,
		},
	}

	return server
}
