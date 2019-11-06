package types

type ContainerStats struct {
	CPUPercentage    float64 `json:"cpu_percentage"`
	MemoryPercentage float64 `json:"memory_percentage"`
	Memory           uint64  `json:"memory"`
	MemoryLimit      uint64  `json:"memory_limit"`
	NetworkDownload  uint64  `json:"network_download"`
	NetworkUpload    uint64  `json:"network_upload"`
	DiscRead         uint64  `json:"disc_read"`
	DiscWrite        uint64  `json:"disc_write"`
}

type DockerContainerConfiguration struct {
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
}

type EventType string
const (
	EventTypeDie EventType = "die"
) 

type ContainerEvent struct {
	ContainerID string
	Event       EventType
}
