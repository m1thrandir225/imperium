package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type ConfigHandler interface {
	GetConfigStatus(ctx *gin.Context)
	SetupConfig(ctx *gin.Context)
}

type configHandler struct {
	config *config.Config

	authService    services.AuthService
	clientService  services.ClientService
	hostService    services.HostService
	sessionService services.SessionService
}

func NewConfigHandler(
	config *config.Config,
) ConfigHandler {
	return &configHandler{
		config: config,
	}
}

func (h *configHandler) GetConfigStatus(ctx *gin.Context) {
	configured := config.IsConfigured()

	ctx.JSON(http.StatusOK, gin.H{
		"configured": configured,
	})
}

func (h *configHandler) SetupConfig(ctx *gin.Context) {
	var req SetupConfigRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if req.AuthServerBaseURL == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("auth server base url is required")))
		return
	}

	if err := config.SaveConfig(req.AuthServerBaseURL); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	h.updateServicesWithNewURL(req.AuthServerBaseURL)

	ctx.JSON(http.StatusOK, messageResponse("Config setup successfully"))
}

// updateServicesWithNewURL updates all services with the new auth server base URL
func (h *configHandler) updateServicesWithNewURL(authServerBaseURL string) {
	apiURL := fmt.Sprintf("%s/api/v1", authServerBaseURL)

	h.authService.UpdateBaseURL(fmt.Sprintf("%s/auth", apiURL))
	h.hostService.UpdateBaseURL(fmt.Sprintf("%s/hosts", apiURL))
	h.clientService.UpdateBaseURL(fmt.Sprintf("%s/clients", apiURL))
	h.sessionService.UpdateBaseURL(fmt.Sprintf("%s/sessions", apiURL))
}
