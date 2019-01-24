package socket

import "github.com/gin-gonic/gin"

type ServerConsolePayload struct {
	ServerId string `json:"server_id"`

	Line string `json:"line"`
}

type ServerStatsUpdatePayload struct {
	ServerId string `json:"server_id"`

	Stats gin.H `json:"stats"`
}

type ServerStatusUpdatePayload struct {
	ServerId string `json:"server_id"`

	Status string `json:"status"`
}
