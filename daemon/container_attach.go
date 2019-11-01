package daemon

import (
	"context"
	"io"

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

	go c.listenEvents()
	go c.attachStats()

	return nil
}
