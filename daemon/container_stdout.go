package daemon

import (
	"fmt"
	"io"
	"strings"

	"code.cloudfoundry.org/bytefmt"
	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/api/socket"
	"github.com/panelmc/daemon/types"
)

var _ io.Writer = &Server{}

// To copy stdout from the docker container, and read directly
func (s *Server) Write(b []byte) (n int, e error) {
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
	payload := socket.ServerConsolePayload{
		ServerID: s.ID,
		Line:     line,
	}
	s.Logger().WithField("event", "Console").Infof("%#v", line)
	socket.BroadcastTo(s.ID, "console_output", payload)

	return len(b), nil
}

func (s *Server) UpdateStats(stats *types.ContainerStats) {
	s.Stats.Usage = *stats
	fStats := gin.H{
		"cpu_percentage":    fmt.Sprintf("%.2f", stats.CPUPercentage),
		"memory_percentage": fmt.Sprintf("%.2f", stats.MemoryPercentage),
		"memory":            bytefmt.ByteSize(stats.Memory),
		"memory_limit":      bytefmt.ByteSize(stats.MemoryLimit),
		"network_download":  bytefmt.ByteSize(stats.NetworkDownload),
		"network_upload":    bytefmt.ByteSize(stats.NetworkUpload),
		"disc_read":         bytefmt.ByteSize(stats.DiscRead),
		"disc_write":        bytefmt.ByteSize(stats.DiscWrite),
	}

	socket.BroadcastTo(s.ID, "stats_update", socket.ServerStatsUpdatePayload{
		ServerID: s.ID,
		Stats:    fStats,
	})
}
