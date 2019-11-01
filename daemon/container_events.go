package daemon

import (
	"context"

	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

func (c *DockerContainer) listenEvents() {
	args := filters.NewArgs()
	args.Add("type", events.ContainerEventType)
	args.Add("container", c.ContainerID)
	args.Add("event", "die")

	msg, err := c.client.Events(context.TODO(), docker.EventsOptions{
		Filters: args,
	})

	for {
		select {
		case message := <-msg:
			if message.ID == c.ContainerID {
				if message.Status == "die" {
					c.server.onDie()
				}
			}
			break
		case erro := <-err:
			logrus.WithField("server", c.server.ID).WithError(erro).Error("An error ocurred on the docker listener.")
			break
		}
	}
}
