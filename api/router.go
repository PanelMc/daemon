package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/api/jwt"
	"github.com/panelmc/daemon/api/routes"
	"github.com/panelmc/daemon/api/socket"
	"github.com/sirupsen/logrus"
)

func Init() {
	router := gin.New()

	//router.Use(gin.Logger())
	// Recover from a panic call :)
	router.Use(gin.Recovery())

	router.LoadHTMLGlob("api/views/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	router.GET("/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"jwt": jwt.NewToken()})
	})
	// Testing view for the socket
	router.GET("/console", func(c *gin.Context) {
		c.HTML(http.StatusOK, "console.tmpl", gin.H{"token": jwt.NewToken()})
	})

	api := router.Group("/api")
	{
		api.Use(jwt.GinHandler)
		api.GET("/", routes.Index)
		api.GET("/servers", routes.ListServers)
		api.GET("/servers/:server", routes.GetServer)
	}

	// socket connection
	if err := socket.Init(); err == nil {
		router.GET("/socket.io/", socket.Handler())
		router.POST("/socket.io/", socket.Handler())
		router.Handle("WS", "/socket.io", socket.Handler())
		router.Handle("WSS", "/socket.io", socket.Handler())
		getLogger().Info("Registered the socket!")
	}

	go func() {
		if err := router.Run(":8080"); err != nil {
			getLogger().WithError(err).Error("Failed to start the daemon API.")
		}

		getLogger().Info("Listening on port 8080!")
	}()
}

func getLogger() *logrus.Entry {
	return logrus.WithField("logger", "API")
}
