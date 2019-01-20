package socket

import (
	"code.cloudfoundry.org/bytefmt"
	"fmt"
)

type ContainerStats struct {
	CPUPercentage    float64
	Memory           float64
	MemoryPercentage float64
	MemoryLimit      float64
	NetworkDownload  float64
	NetworkUpload    float64
	DiscRead         float64
	DiscWrite        float64
}

type formatedContainerStats struct {
	CPUPercentage    string `json:"cpu_percentage"`
	MemoryPercentage string `json:"memory_percentage"`
	Memory           string `json:"memory"`
	MemoryLimit      string `json:"memory_limit"`
	NetworkDownload  string `json:"network_download"`
	NetworkUpload    string `json:"network_upload"`
	DiscRead         string `json:"disc_read"`
	DiscWrite        string `json:"disc_write"`
}

type serverStatsPayload struct {
	ServerId string `json:"server_id"`

	Stats formatedContainerStats `json:"stats"`
}

func (s *ServerRoom) UpdateStats(stats ContainerStats) {
	fStats := formatedContainerStats{
		CPUPercentage:    fmt.Sprintf("%.2f", stats.CPUPercentage),
		MemoryPercentage: fmt.Sprintf("%.2f", stats.MemoryPercentage),
		Memory:           bytefmt.ByteSize(uint64(stats.Memory)),
		MemoryLimit:      bytefmt.ByteSize(uint64(stats.MemoryLimit)),
		NetworkDownload:  bytefmt.ByteSize(uint64(stats.NetworkDownload)),
		NetworkUpload:    bytefmt.ByteSize(uint64(stats.NetworkUpload)),
		DiscRead:         bytefmt.ByteSize(uint64(stats.DiscRead)),
		DiscWrite:        bytefmt.ByteSize(uint64(stats.DiscWrite)),
	}

	Server.BroadcastTo(s.ServerId, "stats_update", serverStatsPayload{s.ServerId, fStats})
}
