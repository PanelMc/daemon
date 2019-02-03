package daemon

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heroslender/panelmc/api/socket"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type ContainerListener struct {
	ServerId string `json:"server_id"`
}

var _ io.Writer = &ContainerListener{}

// To copy stdout from the docker container, and read directly
func (c *ContainerListener) Write(b []byte) (n int, e error) {
	//l := make([]byte, len(b))
	//copy(l, b)
	line := string(b)

	if len(b) == 1 {
		// Is a key pressed, ignore
		return 1, nil
	} else if line == "\r\n" || strings.Contains(line, "\b") {
		return 2, nil
	}

	if strings.Contains(line, "\r") {
		line = strings.Replace(line, "\r", "", -1)
	}
	if strings.Contains(line, "\n") {
		line = strings.Replace(line, "\n", "", -1)
	}

	payload := socket.ServerConsolePayload{c.ServerId, line}
	logrus.WithField("server", c.ServerId).WithField("event", "Console").Infof("%#v", line)
	socket.BroadcastTo(c.ServerId, "console_output", payload)

	return len(b), nil
}

func (c *ContainerListener) UpdateStats(stats ContainerStats) {
	fStats := gin.H{
		"CPUPercentage":    fmt.Sprintf("%.2f", stats.CPUPercentage),
		"MemoryPercentage": fmt.Sprintf("%.2f", stats.MemoryPercentage),
		"Memory":           bytefmt.ByteSize(uint64(stats.Memory)),
		"MemoryLimit":      bytefmt.ByteSize(uint64(stats.MemoryLimit)),
		"NetworkDownload":  bytefmt.ByteSize(uint64(stats.NetworkDownload)),
		"NetworkUpload":    bytefmt.ByteSize(uint64(stats.NetworkUpload)),
		"DiscRead":         bytefmt.ByteSize(uint64(stats.DiscRead)),
		"DiscWrite":        bytefmt.ByteSize(uint64(stats.DiscWrite)),
	}

	socket.BroadcastTo(c.ServerId, "stats_update", socket.ServerStatsUpdatePayload{c.ServerId, fStats})
}
