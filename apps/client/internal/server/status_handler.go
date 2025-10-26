package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatusHandler interface {
	GetStatus(ctx *gin.Context)
}

type statusHandler struct{}

func NewStatusHandler() StatusHandler {
	return &statusHandler{}
}

func (h *statusHandler) GetStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
