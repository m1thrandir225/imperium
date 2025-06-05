package server

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLFiles("web/templates/index.html")
	router.Static("/static", "web/static")

	httpHandler := NewHTTPHandler()

	SetupRoutes(router, httpHandler)

	return router
}
