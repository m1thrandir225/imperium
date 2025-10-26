package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpClient *http.Client

	statusHandler  StatusHandler
	authHandler    AuthHandler
	sessionHandler SessionHandler
	clientHandler  ClientHandler
	hostHandler    HostHandler
	configHandler  ConfigHandler
}

func NewServer(cfg *config.Config) *Server {
	router := newRouter()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiURL := fmt.Sprintf("%s/api/v1", cfg.AuthServerBaseURL)

	authService := services.NewAuthService(fmt.Sprintf("%s/auth", apiURL), httpClient)
	hostService := services.NewHostService(fmt.Sprintf("%s/hosts", apiURL), httpClient)
	clientService := services.NewClientService(fmt.Sprintf("%s/clients", apiURL), httpClient)
	sessionService := services.NewSessionService(fmt.Sprintf("%s/sessions", apiURL), httpClient)

	authHandler, _ := NewAuthHandler(cfg, authService)
	sessionHandler := NewSessionHandler(cfg, sessionService)
	hostHandler, _ := NewHostHandler(cfg, hostService)
	clientHandler, _ := NewClientHandler(cfg, clientService)
	statusHandler := NewStatusHandler()
	configHandler := NewConfigHandler(cfg)

	return &Server{
		config:     cfg,
		router:     router,
		httpClient: httpClient,

		authHandler:    authHandler,
		hostHandler:    hostHandler,
		clientHandler:  clientHandler,
		sessionHandler: sessionHandler,
		statusHandler:  statusHandler,
		configHandler:  configHandler,
	}
}

func newRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.ErrorLogger())

	return router
}

func (s *Server) StartAndServe() {}
