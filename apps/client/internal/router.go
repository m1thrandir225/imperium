package internal

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
)

func SetupRouter(config *config.Config) *gin.Engine {
	router := gin.Default()

	baseURL := config.AuthServerBaseURL

	apiURL := fmt.Sprintf("%s/api/v1", baseURL)

	authService := NewAuthService(apiURL)
	hostService := NewHostService(fmt.Sprintf("%s/hosts", apiURL))
	clientService := NewClientService(fmt.Sprintf("%s/clients", apiURL))
	httpHandler := NewHTTPHandler(config, authService, hostService, clientService)

	router.Use(CORSMiddleware())

	SetupRoutes(router, httpHandler)

	router.StaticFS("/assets", AssetsFS())

	router.GET("/", gin.WrapH(http.FileServer(DistFS())))

	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}
		// Serve index.html for any non-API route
		c.Request.URL.Path = "/"
		router.HandleContext(c)
	})

	return router
}
