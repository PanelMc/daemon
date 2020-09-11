package socket

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/googollee/go-engine.io/transport/websocket"

	"github.com/gin-gonic/gin"
	engineio "github.com/googollee/go-engine.io"
	socketio "github.com/googollee/go-socket.io"
	"github.com/panelmc/daemon/api/jwt"
	"github.com/sirupsen/logrus"
)

var server *socketio.Server

// Init initializes the web socket
func Init() error {
	var err error

	// Fix CORS in go-socket.io
	websocket.Default.CheckOrigin = func(r *http.Request) bool {
		origin := r.Header["Origin"]
		if len(origin) == 0 {
			return true
		}
		u, err := url.Parse(origin[0])
		if err != nil {
			return false
		}

		return strings.EqualFold(u.Host, "localhost:25555") || strings.EqualFold(u.Host, "localhost:8080")
	}

	server, err = socketio.NewServer(&engineio.Options{
		RequestChecker: jwt.SocketHandler,
		ConnInitor: func(r *http.Request, conn engineio.Conn) {
			if serverName := r.FormValue("server"); serverName != "" {
				r.Header.Add("server", serverName)
			}
		},
	})
	if err != nil {
		logrus.WithError(err).Error("can't create socker server.")
		return err
	}

	server.OnConnect("/", func(conn socketio.Conn) error {
		if server := conn.RemoteHeader().Get("server"); server != "" {
			conn.Join(server)
			logrus.Infof("User connected to %s!", conn.RemoteHeader().Get("server"))
		} else {
			logrus.Info("User connected!")
		}

		conn.Join("global")
		return nil
	})

	server.OnEvent("/server", "input", func(s socketio.Conn, server, command string) {

	})
	server.OnError("/", func(err error) {
		logrus.WithError(err).Error("socker server error.")
	})

	go server.Serve()

	return nil
}

// Handler to register on Gin.
func Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		c.Header("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		server.ServeHTTP(c.Writer, c.Request)
	}
}

// Close the socket server
func Close() {
	server.Close()
}

// Broadcast a message to everyone
func Broadcast(event string, messages ...interface{}) {
	BroadcastTo("global", event, messages...)
}

// BroadcastTo a message to everyone in a specific room
func BroadcastTo(room, event string, messages ...interface{}) {
	if server != nil {
		server.BroadcastToRoom(room, event, messages...)
	}
}
