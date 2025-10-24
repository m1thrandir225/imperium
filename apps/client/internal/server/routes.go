package server

func (s *Server) SetupRoutes() {

	apiV1 := s.router.Group("/api/v1")
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

		clients := apiV1.Group("/clients")
		clients.Use(AuthMiddleware())
		{
			clients.GET("", handler.GetClientInfo)
		}

		sessions := apiV1.Group("/sessions")
		sessions.Use(AuthMiddleware())
		{
			sessions.POST("/create", handler.CreateSession)
			sessions.GET("/:sessionId", handler.GetSession)

			sessions.POST("/:sessionId/start", handler.StartSession)
			sessions.POST("/:sessionId/end", handler.EndSession)
			sessions.POST("/:sessionId/cancel", handler.CancelSession)
		}
	}
}
