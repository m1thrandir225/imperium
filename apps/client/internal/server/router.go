package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
	"github.com/m1thrandir225/imperium/apps/client/internal/ui"
)

func SetupRouter(config *config.Config) *gin.Engine {
	router := gin.Default()

	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	router.Use(gin.ErrorLogger())

	baseURL := config.AuthServerBaseURL

	apiURL := fmt.Sprintf("%s/api/v1", baseURL)

	authService := services.NewAuthService(fmt.Sprintf("%s/auth", apiURL))
	hostService := services.NewHostService(fmt.Sprintf("%s/hosts", apiURL))
	clientService := services.NewClientService(fmt.Sprintf("%s/clients", apiURL))
	sessionService := services.NewSessionService(fmt.Sprintf("%s/sessions", apiURL))

	httpHandler := NewHTTPHandler(
		config,
		authService,
		hostService,
		clientService,
		sessionService,
	)

	router.Use(CORSMiddleware())

	SetupRoutes(router, httpHandler)

	router.StaticFS("/assets", ui.AssetsFS())

	router.GET("/", gin.WrapH(http.FileServer(ui.DistFS())))

	router.NoRoute(func(c *gin.Context) {
		if strings.Contains(c.Request.URL.Path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}
		// Serve index.html for any non-API route
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})

	return router
}
