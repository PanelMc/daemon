package api

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/panelmc/daemon/infra"
	"github.com/panelmc/daemon/types"

	"github.com/gin-gonic/gin"
	"github.com/panelmc/daemon/api/jwt"
	"github.com/panelmc/daemon/api/routes"
	"github.com/panelmc/daemon/api/socket"
	"github.com/sirupsen/logrus"
)

func Init(server *infra.Server) {
	router := server.Router
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

	api := router.Group("/api/v1")
	{
		api.BasePath()
		api.Use(jwt.GinHandler)
		api.GET("", routes.Index)
		api.GET("/servers", handle(routes.ListServers))
		api.POST("/servers", handle(routes.CreateServer))
		api.GET("/servers/:server", handle(routes.GetServer))

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

		getLogger().Info("Listening on port 8080...")
		if err := server.Start(); err != nil {
			getLogger().WithError(err).Error("Failed to start the daemon API.")
		}
	}()
}

func handle(handler types.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := handler(c); err != nil {
			switch e := err.(type) {
			case types.APIError:
				c.JSON(int(e.Code), gin.H{
					"error": e,
				})
				break
			default:
				if reflect.Indirect(reflect.ValueOf(e)).NumField() > 1 {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": types.APIError{
							Key:     "server-error",
							Message: err.Error(),
							Extras: types.APIErrorExtras{
								"error": err,
							},
						},
					})
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{
						"error": types.APIError{
							Key:     "server-error",
							Message: err.Error(),
						},
					})
				}
			}
		}
	}
}

func getLogger() *logrus.Entry {
	return logrus.WithField("logger", "API")
}
