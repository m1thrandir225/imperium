package internal

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
)

type HTTPHandler struct {
	config      *config.Config
	authService *AuthService
	hostService *HostService
}

func NewHTTPHandler(config *config.Config, authService *AuthService, hostService *HostService) *HTTPHandler {
	return &HTTPHandler{
		config:      config,
		authService: authService,
		hostService: hostService,
	}
}

func (h *HTTPHandler) GetStatus(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (h *HTTPHandler) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loginResponse, err := h.authService.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, loginResponse)
}

func (h *HTTPHandler) Register(ctx *gin.Context) {
	var req RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	registerUserResponse, err := h.authService.RegisterUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, registerUserResponse)
}

func (h *HTTPHandler) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenResponse, err := h.authService.RefreshToken(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, refreshTokenResponse)
}

func (h *HTTPHandler) GetUserHosts(ctx *gin.Context) {
	authToken := GetAuthToken(ctx)
	if authToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	hosts, err := h.hostService.GetUserHosts(ctx, authToken)
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
