package daemon

import (
	"code.cloudfoundry.org/bytefmt"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/heroslender/panelmc/api/socket"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
	"time"
)

func newDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return cli, err
	}
	cli.NegotiateAPIVersion(context.TODO())
	return cli, nil
}

func containerExists(cli *client.Client, containerId string) bool {
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{All: true})
	if err != nil {
		logrus.WithError(err).Errorf("Failed to check if the container '%s' exists.", containerId)
		return false
	}

	// Check for ID
	for _, c := range containers {
		if c.ID == containerId {
			return true
		}
	}

	// If got none by id, check by name
	for _, c := range containers {
		for _, name := range c.Names {
			if name == containerId {
				return true
			}
		}
	}

	return false
}

// Initialize the docker client and check container
func NewDockerContainer(s *ServerStruct) error {
	cli, err := newDockerClient()
	if err != nil {
		return err
	}
	s.Container.attached = false
	s.Container.server = s
	s.Container.client = cli

	if _, err := cli.ContainerInspect(context.TODO(), s.Container.ContainerId); err != nil {
		// Container wasn't found, setting to an empty string to create a new one later
		s.Container.ContainerId = ""
	}

	s.Save()

	return nil
}

func (c *DockerContainerStruct) Create() error {
	if c.ContainerId != "" {
		logrus.WithField("server", c.server.Id).
			Error("WTF Mate, container already defined /facepalm")
	}

	portSet := nat.PortSet{}
	portMap := nat.PortMap{}
	for _, p := range c.server.Settings.Ports {
		port := nat.Port(fmt.Sprintf("%d/%s", p, "tcp"))
		portSet[port] = struct{}{}
		portMap[port] = []nat.PortBinding{{"0.0.0.0", fmt.Sprintf("%d", p)}}
	}

	containerConfig := &container.Config{
		Image:        "itzg/minecraft-server",
		AttachStdin:  true,
		OpenStdin:    true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
		Hostname:     "daemon-" + c.server.Id,
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

	path := strings.Replace(c.server.DataPath(), "C:\\", "/c/", 1)
	path = strings.Replace(path, "\\", "/", -1)
	path += ":/data"

	memory, err := bytefmt.ToBytes(c.server.Settings.Ram)
	if err != nil {
		logrus.WithField("server", c.server.Id).Error("Failed to read server RAM, using default(1 Gigabyte).")
		memory = 1073741824 // 1GB Default
	}
	swap, err := bytefmt.ToBytes(c.server.Settings.Swap)
	if err != nil {
		logrus.WithField("server", c.server.Id).Error("Failed to read server Swap, using default(1 Gigabyte).")
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

	resContainer, err := c.client.ContainerCreate(context.TODO(), containerConfig, containerHostConfig, nil, containerConfig.Hostname)
	if err != nil {
		return err
	}

	c.ContainerId = resContainer.ID
	c.server.Save()
	return nil
}

func (c *DockerContainerStruct) Attach() error {
	if c.attached {
		return nil
	}

	var err error
	c.hijackedResponse, err = c.client.ContainerAttach(context.TODO(), c.server.Container.ContainerId,
		types.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})

	if err != nil {
		return err
	}
	c.attached = true

	serverRoom := socket.ServerRoom{
		ServerId: c.server.Id,
	}

	go func() {
		defer c.hijackedResponse.Close()
		defer func() {
			c.attached = false
		}()

		if _, err := io.Copy(serverRoom, c.hijackedResponse.Reader); err != nil {
			logrus.WithField("server", c.server.Id).WithError(err).Error("Failed to attach to the server serverRoom!")
		}
	}()

	go func() {
		if !c.attachedStats {
			stats := c.attachStats()
			for {
				serverRoom.UpdateStats(<- stats)
			}
		}
	}()
	return nil
}

func (c *DockerContainerStruct) Start() error {
	logrus.WithField("server", c.server.Id).Debug("Starting the server...")
	if err := c.Attach(); err != nil {
		logrus.WithError(err).Error("Failed to attach to the docker container.")
	}

	if err := c.client.ContainerStart(context.TODO(), c.server.Container.ContainerId, types.ContainerStartOptions{}); err != nil {
		logrus.WithField("server", c.server.Id).Error("Failed to start the docker container.")
		return err
	}
	return nil
}

func (c *DockerContainerStruct) Stop() error {
	timeout := time.Duration(time.Second * 15)

	if err := c.client.ContainerStop(context.TODO(), c.server.Container.ContainerId, &timeout); err != nil {
		logrus.WithField("server", c.server.Id).Error("Failed to stop the docker container.")
		return err
	}
	return nil
}

func (c *DockerContainerStruct) Exec(command string) error {
	_, err := c.hijackedResponse.Conn.Write([]byte(command))
	return err
}
