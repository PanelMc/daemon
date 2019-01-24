package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/heroslender/panelmc/daemon"
	"net/http"
)

func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"user": c.Request.Context().Value("jwt"),
	})
}

func ListServers(c *gin.Context) {
	servers := daemon.GetServers()
	if len(*servers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"error":   "no_servers_loaded",
				"message": "There are no servers loaded.",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": servers,
	})
}

func GetServer(c *gin.Context) {
	server := c.Param("server")

	if s := daemon.GetServer(server); s != nil {
		c.JSON(http.StatusOK, gin.H{
			"data": s,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"error":   "server_not_fount",
				"message": fmt.Sprintf("The server '%s' wasn't found.", server),
			},
		})
	}
}
