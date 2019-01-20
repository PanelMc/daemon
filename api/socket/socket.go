package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/heroslender/panelmc/api/jwt"
	"github.com/sirupsen/logrus"
)

var Server *socketio.Server

func Init() error {
	server, err := socketio.NewServer(nil)
	if err != nil {
		logrus.WithError(err).Error("can't create socker server.")
		return err
	}
	Server = server

	Server.SetAllowRequest(jwt.VerifyRequest)

	Server.On("connection", func(so socketio.Socket) {
		serverName := so.Request().FormValue("server")
		if serverName != "" {
			so.Join(serverName)
			logrus.Infof("User connected to %s!", serverName)
		}

		so.Emit("connected")

		so.On("console_input", func(msg string) {
			// TODO
		})

		so.On("disconnection", func() {
			// TODO
		})
	})

	Server.On("error", func(so socketio.Socket, err error) {
		logrus.WithError(err).Error("socker server error.")
	})

	return nil
}

// Handler initializes the prometheus middleware.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", origin)
		Server.ServeHTTP(c.Writer, c.Request)
	}
}
