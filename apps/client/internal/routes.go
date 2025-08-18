package internal

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/status", handler.GetStatus)
		apiV1.POST("/auth/login", handler.Login)
		apiV1.POST("/auth/register", handler.Register)
		apiV1.POST("/auth/refresh", handler.RefreshToken)
	}
}
