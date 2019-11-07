package daemon

import (
	"context"
	"io"

	"github.com/panelmc/daemon/types"

	docker "github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

// Attach to the given container stdout and live usage stats
func (c *DockerContainer) Attach() error {
	if c.attached {
		return nil
	}

	var err error
	c.hijackedResponse, err = c.client.ContainerAttach(context.TODO(), c.server.Container.ContainerID,
		docker.ContainerAttachOptions{
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Stream: true,
		})

	if err != nil {
		return err
	}

	c.attached = true
	go func() {
		defer func() {
			c.attached = false
			c.hijackedResponse.Close()
		}()

		if _, err := io.Copy(c.server, c.hijackedResponse.Reader); err != nil {
			logrus.WithField("server", c.server.ID).WithError(err).Error("Failed to attach to the server serverRoom!")
		}
	}()

	go func() {
		eventChan := c.listenEvents(context.Background())

		for event := range eventChan {
			if event.ContainerID == c.ContainerID {
				if event.Event == types.EventTypeDie {
					c.server.onDie()
				}
			}
		}
	}()

	go func() {
		statsChan, err := c.attachStats(context.Background())
		if err != nil {
			logrus.WithField("server", c.server.Name).WithError(err).Error("There was an error trying to listen to the container stats")
		}

		for stats := range statsChan {
			c.server.UpdateStats(stats)
		}
	}()

	return nil
}
