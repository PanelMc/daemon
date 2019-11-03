package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/panelmc/daemon/types"

	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/daemon"
)

// Index - Route for GET /
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": c.Request.Context().Value("jwt"),
	})
}

// CreateServer - Route for POST /server
func CreateServer(c *gin.Context) error {
	serverConfig := &types.ServerConfiguration{}
	if err := json.Unmarshal([]byte(c.PostForm("server")), serverConfig); err != nil {
		return types.APIError{
			Code:    http.StatusBadRequest,
			Key:     "server.create.error.invalid-configuration",
			Message: "Failed to parse the configuration, please check your data.",
			Extras: types.APIErrorExtras{
				"error": err,
			},
		}
	}

	server := daemon.NewServer(serverConfig)
	if s := daemon.GetServerByID(server.ID); s != nil {
		return types.APIError{
			Code:    http.StatusBadRequest,
			Key:     "server.create.error.id-in-use",
			Message: fmt.Sprintf("The ID '%s' is already in use.", server.ID),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "Creating the server...",
	})

	go func() {
		if err := server.Init(); err != nil {
			logrus.WithField("context", "API").WithError(err).Error("There was an unexpected error.")
			return
		}

		if err := server.Start(); err != nil {
			logrus.WithField("context", "API").WithError(err).Error("There was an unexpected error.")
		}
	}()

	return nil
}

// ListServers - Route for GET /server
func ListServers(c *gin.Context) error {
	servers := daemon.GetServers()
	if len(*servers) == 0 {
		return types.APIError{
			Code:    http.StatusOK,
			Key:     "server.list.none-avaliable",
			Message: "There are no servers loaded.",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": servers,
	})
	return nil
}

// GetServer - Route for GET /server/{id}
func GetServer(c *gin.Context) error {
	server := c.Param("server")

	if s := daemon.GetServer(server); s != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": s,
		})

		return nil
	}

	return types.APIError{
		Code:    http.StatusNotFound,
		Key:     "server.error.not-found",
		Message: fmt.Sprintf("The server '%s' does not exist.", server),
	}
}
