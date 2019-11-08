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
)

func (c *DockerContainer) attachStats(ctx context.Context, delay time.Duration) (<-chan *types.ContainerStats, error) {
	stats, err := c.client.ContainerStats(ctx, c.ContainerID, true)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		stats.Body.Close()
		return nil, nil
	default:
	}

	statsChan := make(chan *types.ContainerStats)

	go func() {
		defer func() {
			close(statsChan)
			stats.Body.Close()
		}()

		daemonOSType := stats.OSType
		dec := json.NewDecoder(stats.Body)
		var v *docker.StatsJSON

		for {
			if err := dec.Decode(&v); err != nil {
				if err == io.EOF {
					// No more content, exit loop and close everything
					break
				}

				// Create a new decoder with the remaining data from the current decoder
				// in combination with the stats stream reader
				dec = json.NewDecoder(io.MultiReader(dec.Buffered(), stats.Body))

				time.Sleep(delay)
				continue
			}

			statsChan <- mapStats(daemonOSType, v)
		}
	}()

	return statsChan, nil
}

func mapStats(daemonOSType string, v *docker.StatsJSON) *types.ContainerStats {
	var cpuPercent, memPerc float64
	var blkRead, blkWrite, mem, memLimit uint64

	if daemonOSType != "windows" {
		// MemoryStats.Limit will never be 0 unless the container is not running and we haven't
		// got any data from cgroup
		if v.MemoryStats.Limit != 0 {
			memPerc = float64(v.MemoryStats.Usage) / float64(v.MemoryStats.Limit) * 100.0
		}
		cpuPercent = calculateCPUPercentUnix(v.PreCPUStats.CPUUsage.TotalUsage, v.PreCPUStats.SystemUsage, v)
		blkRead, blkWrite = calculateBlockIO(v.BlkioStats)
		mem = v.MemoryStats.Usage
		memLimit = v.MemoryStats.Limit
	} else {
		cpuPercent = calculateCPUPercentWindows(v)
		blkRead = v.StorageStats.ReadSizeBytes
		blkWrite = v.StorageStats.WriteSizeBytes
		mem = v.MemoryStats.PrivateWorkingSet
	}
	netRx, netTx := calculateNetwork(v.Networks)

	return &types.ContainerStats{
		CPUPercentage:    cpuPercent,
		Memory:           mem,
		MemoryPercentage: memPerc,
		MemoryLimit:      memLimit,
		NetworkDownload:  netRx,
		NetworkUpload:    netTx,
		DiscRead:         blkRead,
		DiscWrite:        blkWrite,
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

func calculateNetwork(network map[string]docker.NetworkStats) (uint64, uint64) {
	var rx, tx uint64

	for _, v := range network {
		rx += v.RxBytes
		tx += v.TxBytes
	}

	return rx, tx
}
