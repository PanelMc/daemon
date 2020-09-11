package types

type Player struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

type ServerStats struct {
	Status        ServerStatus   `json:"status"`
	OnlinePlayers []*Player      `json:"online_players"`
	Usage         ContainerStats `json:"usage"`
}

type ServerSettings struct {
	Ram   string `json:"ram"`
	Swap  string `json:"swap"`
	Ports []int  `json:"ports"`
}

type ServerStatus string

func (s ServerStatus) String() string {
	return string(s)
}

const (
	ServerStatusOnline   ServerStatus = "online"
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusOffline  ServerStatus = "offline"
	ServerStatusStopping ServerStatus = "stopping"
)

type ServerConfiguration struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`

	Settings ServerSettings `json:"settings"`

	Container DockerContainerConfiguration `json:"container"`
}