package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/api/jwt"
	"github.com/panelmc/daemon/api/routes"
	"github.com/panelmc/daemon/api/socket"
	"github.com/sirupsen/logrus"
)

func Init() {
	router := gin.New()
	router.RedirectTrailingSlash = true

	router.Use(gin.Logger())
	// Recover from a panic call :)
	var noRouteHandlers []gin.HandlerFunc
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

	// Redirect /api to /api/v1
	router.Any("/api", func(c *gin.Context) {
		c.Request.URL.Path = "/api/v1"
		router.HandleContext(c)
	})

	api := router.Group("/api/v1")
	{
		api.BasePath()
		api.Use(jwt.GinHandler)
		api.GET("", routes.Index)
		api.GET("/servers", routes.ListServers)
		api.GET("/servers/:server", routes.GetServer)

		noRouteHandlers = append(noRouteHandlers, func(c *gin.Context) {
			if strings.HasPrefix(c.Request.URL.String(), api.BasePath()) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": gin.H{
						"error":   "404",
						"message": "API endpoint not found.",
					},
				})
			}
		})

		getLogger().Info("Registered the api route!")
	}

	// socket connection
	if err := socket.Init(); err == nil {
		router.GET("/socket.io/", socket.Handler())
		router.POST("/socket.io/", socket.Handler())
		router.Handle("WS", "/socket.io", socket.Handler())
		router.Handle("WSS", "/socket.io", socket.Handler())
		getLogger().Info("Registered the socket!")
	}

	router.NoRoute(noRouteHandlers...)

	go func() {
		defer socket.Close()

		if err := router.Run(":8080"); err != nil {
			getLogger().WithError(err).Error("Failed to start the daemon API.")
		}

		getLogger().Info("Listening on port 8080!")
	}()
}

func getLogger() *logrus.Entry {
	return logrus.WithField("logger", "API")
}
