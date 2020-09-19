package infra

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// ServerMode indicates in wich mode the server is running
type ServerMode string

func (s *ServerMode) String() string {
	return string(*s)
}

const (
	// DebugMode indicates gin mode is debug.
	DebugMode ServerMode = gin.DebugMode
	// ReleaseMode indicates gin mode is release.
	ReleaseMode ServerMode = gin.ReleaseMode
	// TestMode indicates gin mode is test.
	TestMode ServerMode = gin.TestMode
)

// Server is the webserver running the API
type Server struct {	
	Port   uint16
	Router *gin.Engine
}

// NewServer initializes a new server instance for the API
func NewServer(port uint16, mode ServerMode) *Server {
	server := &Server{
		Port:   port,
		Router: gin.New(),
	}

	gin.SetMode(mode.String())

	server.SetCors("*")

	server.Router.RedirectTrailingSlash = true

	return server
}

// SetCors is a helper to set current engine cors
func (s *Server) SetCors(allowedOrigins string) {
	s.Router.Use(cors.Middleware(cors.Config{
		Origins:         allowedOrigins,
		Methods:         strings.Join([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodOptions, http.MethodPatch}, ","),
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
}

// Start starts the server on the configured port
func (s *Server) Start() (err error) {
	err = s.Router.Run(":" + strconv.Itoa(int(s.Port)))
	return
}
