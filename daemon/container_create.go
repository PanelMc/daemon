package daemon

import (
	"context"
	"fmt"
	"strings"

	"github.com/panelmc/daemon/types"

	"code.cloudfoundry.org/bytefmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

// Create a new Docker Container based on this configuration
func (c *DockerContainer) Create() error {
	if c.ContainerID != "" {
		return types.APIError{
			Code:    400,
			Key:     "container.create.error.already_defined",
			Message: "WTF Mate, container already defined /facepalm",
			Extras: &types.APIErrorExtras{
				"server_id": c.server.ID,
			},
		}
	}

	ctx := context.Background()

	containerConfig := parseContainerConfig(c)
	containerHostConfig := parseHostConfig(c)

	resContainer, err := c.client.ContainerCreate(ctx, containerConfig, containerHostConfig, nil, containerConfig.Hostname)
	if err != nil {
		return err
	}

	c.ContainerID = resContainer.ID
	c.server.Save()
	return nil
}

func parseContainerConfig(c *DockerContainer) *container.Config {
	portSet := nat.PortSet{}
	for _, p := range c.server.Settings.Ports {
		port := nat.Port(fmt.Sprintf("%d/%s", p, "tcp"))
		portSet[port] = struct{}{}
	}

	containerConfig := &container.Config{
		Image:        c.Image,
		AttachStdin:  true,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Hostname:     "daemon-" + c.server.ID,
		ExposedPorts: portSet,
		Volumes: map[string]struct{}{
			"/data": {},
		},
		// TODO set ram on env variables
		Env: []string{
			"EULA=TRUE",
			"PAPER_DOWNLOAD_URL=https://heroslender.com/assets/PaperSpigot-1.8.8.jar",
			"TYPE=PAPER",
			"VERSION=1.8.8",
			"ENABLE_RCON=false",
		},
	}

	return containerConfig
}

func parseHostConfig(c *DockerContainer) *container.HostConfig {
	portMap := nat.PortMap{}
	for _, p := range c.server.Settings.Ports {
		port := nat.Port(fmt.Sprintf("%d/%s", p, "tcp"))
		portMap[port] = []nat.PortBinding{{"0.0.0.0", fmt.Sprintf("%d", p)}}
	}

	// fix windows path
	path := strings.Replace(c.server.DataPath(), "C:\\", "/c/", 1)
	path = strings.Replace(path, "\\", "/", -1)
	// point to `/data` volume
	path += ":/data"

	memory, err := bytefmt.ToBytes(c.server.Settings.Ram)
	if err != nil {
		c.server.Logger().Error("Failed to read server RAM, using default(1 Gigabyte).")
		memory = 1073741824 // 1GB Default
	}
	swap, err := bytefmt.ToBytes(c.server.Settings.Swap)
	if err != nil {
		c.server.Logger().Error("Failed to read server Swap, using default(1 Gigabyte).")
		swap = 1073741824 // 1GB Default
	}
	containerHostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:     int64(memory),
			MemorySwap: int64(swap),
		},
		Binds:        []string{path},
		PortBindings: portMap,
	}

	return containerHostConfig
}
