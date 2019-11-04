package daemon

import (
	"context"

	docker "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/panelmc/daemon/types"
)

type IDockerContainer interface {
	Create() error

	Attach() error
	Start() error
	Stop() error

	Exec(command string) error
}

type DockerContainer struct {
	ContainerID string `json:"container_id" jsonapi:"attr,container_id"`
	Image       string `json:"image"`

	client *client.Client

	attachedStats    bool
	attached         bool
	hijackedResponse docker.HijackedResponse

	server *Server
}

var _ IDockerContainer = &DockerContainer{}

// NewDockerContainer - Initialize the docker client and check container
func NewDockerContainer(s *Server) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	s.Container.attached = false
	s.Container.server = s
	s.Container.client = cli
	s.Stats = &types.ServerStats{
		Status:        types.ServerStatusOffline,
		OnlinePlayers: []*types.Player{},
	}

	// Container already has an ID, check if it exists
	if s.Container.ContainerID != "" {
		if _, err := cli.ContainerInspect(context.TODO(), s.Container.ContainerID); err != nil {
			// Container wasn't found, setting to an empty string to create a new one later
			s.Container.ContainerID = ""
		}
	}

	s.Save()

	return nil
}

// Prepare the container, ensures image is pulled and creates the container
func (c *DockerContainer) prepare(ctx context.Context) error {
	if c.ContainerID != "" {
		// Container already created?
		if _, err := c.client.ContainerInspect(context.TODO(), c.ContainerID); err != nil {
			// Container does not exist after all
			return err
		}

		// It is created already, just return nil
		return nil
	}

	if err := c.pullImage(ctx); err != nil {
		return err
	}

	if err := c.Create(); err != nil {
		return err
	}

	return nil
}
