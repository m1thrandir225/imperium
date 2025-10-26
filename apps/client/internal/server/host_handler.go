package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type HostHandler interface {
	GetUserHosts(ctx *gin.Context)
	GetHost(ctx *gin.Context)
	GetHostPrograms(ctx *gin.Context)
}

type hostHandler struct {
	config  *config.Config
	service services.HostService
}

func NewHostHandler(
	cfg *config.Config,
	service services.HostService,
) (HostHandler, error) {
	return &hostHandler{
		config:  cfg,
		service: service,
	}, nil
}

func (h *hostHandler) GetUserHosts(ctx *gin.Context) {
	authToken := GetAuthToken(ctx)

	if authToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hosts, err := h.service.GetUserHosts(ctx, authToken)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, hosts)
}
func (h *hostHandler) GetHost(ctx *gin.Context) {
	token := GetAuthToken(ctx)

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	var uriID HostUriID
	if err := ctx.ShouldBindUri(&uriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	host, err := h.service.GetHost(ctx, uriID.HostID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, host)
}
func (h *hostHandler) GetHostPrograms(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	var uriID HostUriID
	if err := ctx.ShouldBindUri(&uriID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	programs, err := h.service.GetHostPrograms(ctx, uriID.HostID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, programs)
}
