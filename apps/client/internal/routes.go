package internal

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/status", handler.GetStatus)

		config := apiV1.Group("/config")
		{
			config.POST("/setup", handler.SetupConfig)
			config.GET("/status", handler.GetConfigStatus)
		}

		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/register", handler.Register)
			auth.POST("/refresh-token", handler.RefreshToken)
		}
		hosts := apiV1.Group("/hosts")
		hosts.Use(AuthMiddleware())
		{
			hosts.GET("", handler.GetUserHosts)
			hosts.GET("/:hostId", handler.GetHost)
			hosts.GET("/:hostId/programs", handler.GetHostPrograms)
		}
	}
}
