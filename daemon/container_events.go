package daemon

import (
	"context"

	"github.com/panelmc/daemon/types"

	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
)

// TODO - Global event listener instead of a per-server listener
func (c *DockerContainer) listenEvents(ctx context.Context) <-chan types.ContainerEvent {
	args := filters.NewArgs()
	args.Add("type", events.ContainerEventType)
	args.Add("container", c.ContainerID)
	args.Add("event", "die")

	msg, err := c.client.Events(ctx, docker.EventsOptions{
		Filters: args,
	})

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	eventChannel := make(chan types.ContainerEvent)

	go func() {
		defer func() {
			close(eventChannel)
		}()

	loop:
		for {
			select {
			case message := <-msg:
				var eventType types.EventType
				switch message.Status {
				case "die":
					eventType = types.EventTypeDie
					break
				default:
					eventType = types.EventType(message.Status)
				}

				eventChannel <- types.ContainerEvent{
					ContainerID: message.Actor.ID,
					Event:       eventType,
				}
				break
			case <-ctx.Done():
			case <-err:
				break loop
			}
		}
	}()

	return eventChannel
}
