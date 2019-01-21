package api

import (
	"fmt"
	jwt2 "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/heroslender/panelmc/api/jwt"
	"github.com/heroslender/panelmc/daemon"
	"net/http"
)

func initApiRouter(api *gin.RouterGroup) {
	api.Use(func(c *gin.Context) {
		if err := jwt.VerifyRequest(c.Request); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": gin.H{
					"error":   "unauthorized",
					"message": "You need to login to access this content.",
				},
			})
		}
	})
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user": c.Request.Context().Value("jwt").(*jwt2.Token).Claims,
		})
	})
	api.GET("/servers", listServers)
	api.GET("/servers/:server", getServer)
}

func listServers(c *gin.Context) {
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

func getServer(c *gin.Context) {
	server := c.Param("server")
	servers := daemon.GetServers()

	for _, s := range *servers {
		if s.Id == server || s.Name == server {
			c.JSON(http.StatusOK, gin.H{
				"data": s,
			})
			return
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"error": gin.H{
			"error":   "server_not_fount",
			"message": fmt.Sprintf("The server '%s' wasn't found.", server),
		},
	})
}
