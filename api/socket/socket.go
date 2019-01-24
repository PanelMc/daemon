package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"github.com/heroslender/panelmc/api/jwt"
	"github.com/sirupsen/logrus"
)

var socket *socketio.Server

func Init() error {
	var err error
	socket, err = socketio.NewServer(nil)
	if err != nil {
		logrus.WithError(err).Error("can't create socker server.")
		return err
	}

	socket.SetAllowRequest(jwt.VerifyRequest)
	socket.On("connection", func(so socketio.Socket) {
		serverName := so.Request().FormValue("server")
		if serverName != "" {
			so.Join(serverName)
			logrus.Infof("User connected to %s!", serverName)
		}

		so.Join("global")
		so.Emit("connected")

		so.On("console_input", func(msg string) {
			// TODO
		})

		so.On("disconnection", func() {
			// TODO
		})
	})

	socket.On("error", func(so socketio.Socket, err error) {
		logrus.WithError(err).Error("socker server error.")
	})

	return nil
}

// Handler to register on Gin.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", origin)
		socket.ServeHTTP(c.Writer, c.Request)
	}
}

// Broadcast a message to everyone
func Broadcast(event string, messages ...interface{}) {
	BroadcastTo("global", event, messages)
}

// Broadcast a message to everyone in a specific room
func BroadcastTo(room, event string, messages ...interface{}) {
	if socket != nil {
		socket.BroadcastTo(room, event, messages)
	}
}
