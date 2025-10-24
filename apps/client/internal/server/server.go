package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
)

type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpClient *http.Client

	authHandler    AuthHandler
	sessionHandler SessionHandler
	clientHandler  ClientHandler
	hostHandler    HostHandler
}

func NewServer() *Server {
	return &Server{}
}
