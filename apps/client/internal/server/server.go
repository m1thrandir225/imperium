package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
	"github.com/m1thrandir225/imperium/apps/client/internal/ui"
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

// NewServer returns a new server instance
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

	authHandler, _ := NewAuthHandler(authService, clientService)
	sessionHandler := NewSessionHandler(sessionService)
	hostHandler, _ := NewHostHandler(hostService)
	clientHandler, _ := NewClientHandler(clientService)
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

// newRouter creates a gin.Engine and configures it to use default middleware
func newRouter() *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.ErrorLogger())
	router.Use(CORSMiddleware())

	return router
}

func (s *Server) SetupRouter() error {
	s.setupRoutes()

	s.router.StaticFS("/assets", ui.AssetsFS())
	s.router.GET("/", gin.WrapH(http.FileServer(ui.DistFS())))

	s.router.NoRoute(func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.Path, "/api/") {
			ctx.Status(http.StatusNotFound)
			return
		}
		// Serve index.html for any non-API route
		ctx.Request.URL.Path = "/"
		s.router.HandleContext(ctx)
	})

	return nil
}

// StartAndServe
func (s *Server) StartAndServe() {
	server := &http.Server{
		Addr:    ":8081",
		Handler: s.router.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
