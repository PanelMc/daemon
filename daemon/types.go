package daemon

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Server interface {
	Start() error
	Stop() error
	Restart() error

	Execute(command string) error

	Init() error
	DataPath() string
	Save() error
}

type ServerStruct struct {
	Id   string `json:"id"`
	Name string `json:"name"`

	Settings ServerSettings `json:"settings"`

	Stats *ServerStats `json:"stats,omitempty"`

	Container DockerContainerStruct `json:"container"`
}

type ServerSettings struct {
	Ram   string `json:"ram"`
	Swap  string `json:"swap"`
	Ports []int  `json:"ports"`
}

var _ Server = &ServerStruct{}

type DockerContainer interface {
	Create() error

	Attach() error
	Start() error
	Stop() error

	Exec(command string) error
}

type DockerContainerStruct struct {
	ContainerId string `json:"container_id" jsonapi:"attr,container_id"`
	Image       string `json:"image"`

	client *client.Client

	attachedStats    bool
	attached         bool
	hijackedResponse types.HijackedResponse

	server *ServerStruct
}

var _ DockerContainer = &DockerContainerStruct{}

type ServerStats struct {
	Status ServerStatus   `json:"status"`
	Usage  ContainerStats `json:"usage"`
}

type ContainerStats struct {
	CPUPercentage    float64
	Memory           float64
	MemoryPercentage float64
	MemoryLimit      float64
	NetworkDownload  float64
	NetworkUpload    float64
	DiscRead         float64
	DiscWrite        float64
}

type ServerStatus string

const (
	ServerStatusOnline   ServerStatus = "online"
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusOffline  ServerStatus = "offline"
	ServerStatusStopping ServerStatus = "stopping"
)

type ApiError struct {
	Err   string `json:"error"`
	Message string `json:"message"`
}

func (e ApiError) Error() string {
	return e.Message
}
