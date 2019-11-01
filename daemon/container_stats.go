// Code from Moby - https://github.com/moby/moby/blob/eb131c5383db8cac633919f82abad86c99bffbe5/cli/command/container/stats_helpers.go

package daemon

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"time"

	docker "github.com/docker/docker/api/types"
	"github.com/panelmc/daemon/types"
	"github.com/sirupsen/logrus"
)

func (c *DockerContainer) attachStats() {
	c.attachedStats = true
	stats, err := c.client.ContainerStats(context.TODO(), c.ContainerID, true)
	if err != nil {
		logrus.WithField("server", c.server.ID).WithError(err).Error("Failed to read docker container stats.")
		c.attachedStats = false
	}

	defer func() {
		stats.Body.Close()
		c.attachedStats = false
	}()

	daemonOSType := stats.OSType
	dec := json.NewDecoder(stats.Body)

	for {
		var v *docker.StatsJSON

		if err := dec.Decode(&v); err != nil {
			if err == io.EOF {
				// No more content
				go func() {
					time.Sleep(100 * time.Millisecond)
					if !c.attachedStats {
						c.attachStats()
					}
				}()
				break
			}

			dec = json.NewDecoder(io.MultiReader(dec.Buffered(), stats.Body))
			time.Sleep(100 * time.Millisecond)
			continue
		}

		c.server.UpdateStats(*mapStats(daemonOSType, v))
	}
}

func mapStats(daemonOSType string, v *docker.StatsJSON) *types.ContainerStats {
	var (
		memPercent        = 0.0
		cpuPercent        = 0.0
		blkRead, blkWrite uint64 // Only used on Linux
		mem               = 0.0
		memLimit          = 0.0
		memPerc           = 0.0
	)

	if daemonOSType != "windows" {
		// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
		// got any data from cgroup
		if v.MemoryStats.Limit != 0 {
			memPercent = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
		}
		cpuPercent = calculateCPUPercentUnix(v.PreCPUStats.CPUUsage.TotalUsage, v.PreCPUStats.SystemUsage, v)
		blkRead, blkWrite = calculateBlockIO(v.BlkioStats)
		mem = float64(v.MemoryStats.Usage)
		memLimit = float64(v.MemoryStats.Limit)
		memPerc = memPercent
	} else {
		cpuPercent = calculateCPUPercentWindows(v)
		blkRead = v.StorageStats.ReadSizeBytes
		blkWrite = v.StorageStats.WriteSizeBytes
		mem = float64(v.MemoryStats.PrivateWorkingSet)
	}
	netRx, netTx := calculateNetwork(v.Networks)

	return &types.ContainerStats{
		CPUPercentage:    cpuPercent,
		Memory:           mem,
		MemoryPercentage: memPerc,
		MemoryLimit:      memLimit,
		NetworkDownload:  netRx,
		NetworkUpload:    netTx,
		DiscRead:         float64(blkRead),
		DiscWrite:        float64(blkWrite),
	}
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *docker.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)

	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * float64(len(v.CPUStats.CPUUsage.PercpuUsage)) * 100.0
	}
	return cpuPercent
}

func calculateCPUPercentWindows(v *docker.StatsJSON) float64 {
	// Max number of 100ns intervals between the previous time read and now
	possIntervals := uint64(v.Read.Sub(v.PreRead).Nanoseconds()) // Start with number of ns intervals
	possIntervals /= 100                                         // Convert to number of 100ns intervals
	possIntervals *= uint64(v.NumProcs)                          // Multiple by the number of processors

	// Intervals used
	intervalsUsed := v.CPUStats.CPUUsage.TotalUsage - v.PreCPUStats.CPUUsage.TotalUsage

	// Percentage avoiding divide-by-zero
	if possIntervals > 0 {
		return float64(intervalsUsed) / float64(possIntervals) * 100.0
	}
	return 0.00
}

func calculateBlockIO(blkio docker.BlkioStats) (blkRead uint64, blkWrite uint64) {
	for _, bioEntry := range blkio.IoServiceBytesRecursive {
		switch strings.ToLower(bioEntry.Op) {
		case "read":
			blkRead = blkRead + bioEntry.Value
		case "write":
			blkWrite = blkWrite + bioEntry.Value
		}
	}
	return
}

func calculateNetwork(network map[string]docker.NetworkStats) (float64, float64) {
	var rx, tx float64

	for _, v := range network {
		rx += float64(v.RxBytes)
		tx += float64(v.TxBytes)
	}
	return rx, tx
}
