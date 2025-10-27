package server

func (s *Server) setupRoutes() {
	apiV1 := s.router.Group("/api/v1")
	{
		apiV1.GET("/status", s.statusHandler.GetStatus)

		config := apiV1.Group("/config")
		{
			config.POST("/setup", s.configHandler.SetupConfig)
			config.GET("/status", s.configHandler.GetConfigStatus)
		}

		auth := apiV1.Group("/auth")
		{
			auth.POST("/login", s.authHandler.Login)
			auth.POST("/register", s.authHandler.Register)
			auth.POST("/refresh-token", s.authHandler.RefreshToken)
		}
		hosts := apiV1.Group("/hosts")
		hosts.Use(AuthMiddleware())
		{
			hosts.GET("", s.hostHandler.GetUserHosts)
			hosts.GET("/:hostId", s.hostHandler.GetHost)
			hosts.GET("/:hostId/programs", s.hostHandler.GetHostPrograms)
		}

		clients := apiV1.Group("/clients")
		clients.Use(AuthMiddleware())
		{
			clients.GET("", s.clientHandler.GetClientInfo)
		}

		sessions := apiV1.Group("/sessions")
		sessions.Use(AuthMiddleware())
		{
			sessions.POST("/create", s.sessionHandler.CreateSession)
			sessions.GET("/:sessionId", s.sessionHandler.GetSession)

			sessions.POST("/:sessionId/start", s.sessionHandler.StartSession)
			sessions.POST("/:sessionId/end", s.sessionHandler.EndSession)
			sessions.POST("/:sessionId/cancel", s.sessionHandler.CancelSession)
		}
	}
}
