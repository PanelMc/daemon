package daemon

import (
	"strings"
	"testing"

	"github.com/panelmc/daemon/types"
)

func mockDockerContainer() *DockerContainer {
	return &DockerContainer{
		Image: "itzg/minecraft-server",
		server: &Server{
			ID:   "test",
			Name: "Test Server",
			Settings: types.ServerSettings{
				Ram:   "2GB",
				Swap:  "3GB",
				Ports: []int{25565},
			},
			Stats: &types.ServerStats{
				Status:        types.ServerStatusOffline,
				OnlinePlayers: []*types.Player{},
			},
		},
	}
}

func TestContainerConfigurationParser(t *testing.T) {
	c := mockDockerContainer()

	containerConfig := parseContainerConfig(c)
	if containerConfig.Image != c.Image {
		t.Errorf("Image parse failed! Got %s, expected %s.", containerConfig.Image, c.Image)
	}
	if containerConfig.Hostname != "daemon-"+c.server.ID {
		t.Errorf("Container name invalid! Got %s, expected %s", containerConfig.Hostname, "daemon-"+c.server.ID)
	}
}

func TestHostConfigurationParser(t *testing.T) {
	c := mockDockerContainer()

	hostConfig := parseHostConfig(c)
	if hostConfig.Resources.Memory != 2147483648 {
		t.Errorf("Memory parse failed! Got %d, expected %d.", hostConfig.Resources.Memory, 2147483648)
	}
	if hostConfig.Resources.MemorySwap != 3221225472 {
		t.Errorf("Memory Swap parse failed! Got %d, expected %d.", hostConfig.Resources.MemorySwap, 3221225472)
	}
	if !strings.HasSuffix(hostConfig.Binds[0], ":/data") {
		t.Errorf("Volume bind parse failed! Got %s, expected to end with %s.", hostConfig.Binds[0], ":/data")
	}
}
