package daemon

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/panelmc/daemon/types"

	docker "github.com/docker/docker/api/types"
)

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

var pendingImagePulls = []string{}

func (c *DockerContainer) pullImage(ctx context.Context) error {
	for _, img := range pendingImagePulls {
		if strings.EqualFold(img, c.Image) {
			return types.APIError{
				Code:    http.StatusOK,
				Key:     "image.pull.error.pending",
				Message: "Image is already being pulled, please wait.",
			}
		}
	}
	pendingImagePulls = append(pendingImagePulls, c.Image)

	defer func() {
		for i, img := range pendingImagePulls {
			if strings.EqualFold(img, c.Image) {
				pendingImagePulls = append(pendingImagePulls[:i], pendingImagePulls[i+1:]...)
			}
		}
	}()

	logrus.WithField("context", "Daemon").Infof("Pulling image %s...", c.Image)
	r, err := c.client.ImagePull(ctx, c.Image, docker.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer r.Close()

	d := json.NewDecoder(r)
	var event *Event
	for {
		if err := d.Decode(&event); err != nil && err == io.EOF {
			break
		}
	}

	// Return error from last message if present
	if event.Error != "" {
		return errors.New(event.Error)
	}

	logrus.WithField("context", "Daemon").Infof("Image %s pulled!", c.Image)
	return nil
}
