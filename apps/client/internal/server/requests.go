package server

import (
	"github.com/gin-gonic/gin"
)

type HostUriID struct {
	HostID string `uri:"hostId" binding:"required,uuid"`
}

type SetupConfigRequest struct {
	AuthServerBaseURL string `json:"auth_server_base_url",binding:"required"`
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func messageResponse(message string) gin.H {
	return gin.H{"message": message}
}
