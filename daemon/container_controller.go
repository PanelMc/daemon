package daemon

import (
	"context"
	"fmt"
	"time"

	docker "github.com/docker/docker/api/types"
	"github.com/panelmc/daemon/types"
	"github.com/pkg/errors"
)

func (c *DockerContainer) Start() error {
	if c.server.Stats.Status != types.ServerStatusOffline {
		return errors.New(fmt.Sprintf("Server already running. Current status: %s", c.server.Stats.Status))
	}

	c.server.Logger().Debug("Starting the server...")
	if err := c.Attach(); err != nil {
		c.server.Logger().WithError(err).Error("Failed to attach to the docker container.")
	}

	if err := c.client.ContainerStart(context.TODO(), c.server.Container.ContainerID, docker.ContainerStartOptions{}); err != nil {
		c.server.Logger().Error("Failed to start the docker container.")
		return err
	}

	return nil
}

func (c *DockerContainer) Stop() error {
	if c.server.Stats.Status != types.ServerStatusOnline {
		return types.APIError{
			Code:    400,
			Key:     "container.stop.error.already_offline",
			Message: fmt.Sprintf("Server isn't running. Current status: %s", c.server.Stats.Status),
		}
	}

	timeout := time.Duration(time.Second * 15)

	if err := c.client.ContainerStop(context.TODO(), c.server.Container.ContainerID, &timeout); err != nil {
		c.server.Logger().Error("Failed to stop the docker container.")
		return err
	}

	return nil
}

func (c *DockerContainer) Exec(command string) error {
	_, err := c.hijackedResponse.Conn.Write([]byte(command))
	return err
}
