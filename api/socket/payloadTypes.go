package socket

import (
	"github.com/gin-gonic/gin"
)

// ServerConsolePayload - Payload for a console logs
type ServerConsolePayload struct {
	ServerID string `json:"server_id"`

	Line string `json:"line"`
}

// ServerStatsUpdatePayload - Payload for usage stats
type ServerStatsUpdatePayload struct {
	ServerID string `json:"server_id"`

	Stats gin.H `json:"stats"`
}

// ServerStatusUpdatePayload - Payload for status update
type ServerStatusUpdatePayload struct {
	ServerID string `json:"server_id"`

	Status string `json:"status"`
}

// ServerPlayerJoinPayload - Payload for a new player join
type ServerPlayerJoinPayload struct {
	ServerID string `json:"server_id"`

	Player interface{} `json:"player"`
}

// ServerPlayerLeavePayload - Payload for a player disconnect
type ServerPlayerLeavePayload struct {
	ServerID string `json:"server_id"`

	Player interface{} `json:"player"`
}
