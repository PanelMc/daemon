package types

type ContainerStats struct {
	CPUPercentage    float64 `json:"cpu_percentage"`
	Memory           float64 `json:"memory"`
	MemoryPercentage float64 `json:"memory_percentage"`
	MemoryLimit      float64 `json:"memory_limit"`
	NetworkDownload  float64 `json:"network_download"`
	NetworkUpload    float64 `json:"network_upload"`
	DiscRead         float64 `json:"disc_read"`
	DiscWrite        float64 `json:"disc_write"`
}

type DockerContainerConfiguration struct {
	ContainerID string `json:"container_id"`
	Image       string `json:"image"`
}
