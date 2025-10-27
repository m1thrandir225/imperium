package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/internal/models"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type ClientHandler interface {
	GetClientInfo(ctx *gin.Context)
}

type clientHandler struct {
	service services.ClientService
}

func NewClientHandler(
	service services.ClientService,
) (ClientHandler, error) {
	return &clientHandler{
		service: service,
	}, nil
}

func (h *clientHandler) GetClientInfo(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	hostname, ip, err := h.service.GetClientInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientReq := models.RegisterClientRequest{
		Name:      hostname,
		IPAddress: ip,
	}

	clientResponse, err := h.service.RegisterOrUpdateClient(ctx, clientReq, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, clientResponse)
}
