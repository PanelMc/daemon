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
	Id   string `json:"id" jsonapi:"primary,id"`
	Name string `json:"name" jsonapi:"attr,name"`

	Settings ServerSettings `json:"settings" jsonapi:"relation,settings"`

	Container DockerContainerStruct `json:"container" jsonapi:"relation,container"`
}

type ServerSettings struct {
	Ram   string `json:"ram" jsonapi:"attr,ram"`
	Swap  string `json:"swap" jsonapi:"attr,ram"`
	Ports []int  `json:"ports" jsonapi:"attr,name"`
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
