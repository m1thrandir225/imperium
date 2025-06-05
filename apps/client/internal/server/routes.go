package server

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, handler *HTTPHandler) {
	routes := router.Group("/")
	{
		routes.GET("/", handler.ServeIndex)

		apiV1 := routes.Group("/api/v1")
		{
			apiV1.GET("/status", handler.GetStatus)
		}
	}
}
