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

var _ io.Writer = &ServerStruct{}

// To copy stdout from the docker container, and read directly
func (s *ServerStruct) Write(b []byte) (n int, e error) {
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

	processConsoleOutput(s, line)
	payload := socket.ServerConsolePayload{s.Id, line}
	logrus.WithField("server", s.Id).WithField("event", "Console").Infof("%#v", line)
	socket.BroadcastTo(s.Id, "console_output", payload)

	return len(b), nil
}

func (s *ServerStruct) UpdateStats(stats ContainerStats) {
	s.Stats.Usage = stats
	fStats := gin.H{
		"cpu_percentage":    fmt.Sprintf("%.2f", stats.CPUPercentage),
		"memory_percentage": fmt.Sprintf("%.2f", stats.MemoryPercentage),
		"memory":           bytefmt.ByteSize(uint64(stats.Memory)),
		"memory_limit":      bytefmt.ByteSize(uint64(stats.MemoryLimit)),
		"network_download":  bytefmt.ByteSize(uint64(stats.NetworkDownload)),
		"network_upload":    bytefmt.ByteSize(uint64(stats.NetworkUpload)),
		"disc_read":         bytefmt.ByteSize(uint64(stats.DiscRead)),
		"disc_write":        bytefmt.ByteSize(uint64(stats.DiscWrite)),
	}

	socket.BroadcastTo(s.Id, "stats_update", socket.ServerStatsUpdatePayload{s.Id, fStats})
}
