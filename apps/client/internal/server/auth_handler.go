package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/models"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
	"github.com/m1thrandir225/imperium/apps/client/internal/util"
)

type AuthHandler interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type authHandler struct {
	config  *config.Config
	service services.AuthService
}

func NewAuthHandler(
	cfg *config.Config,
	service services.AuthService,
) (AuthHandler, error) {
	return &authHandler{
		config:  cfg,
		service: service,
	}, nil
}

func (h *authHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	loginResponse, err := h.service.Login(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	channelClient := make(chan models.Client)

	// upsert client
	go func() {
		hostname, err := util.GetHostname()
		if err != nil {
			hostname = "unknown-hostname"
		}
		ip, err := util.GetIPV4Address()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		clientResponse, err := h.clientService.RegisterOrUpdateClient(ctx, models.RegisterClientRequest{
			Name:      hostname,
			IPAddress: ip,
		}, loginResponse.AccessToken)
		if err != nil {
			log.Println("error upserting client", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		} else {
			log.Println("client upserted successfully")
			client := models.Client{
				ID:        clientResponse.ID,
				Name:      clientResponse.Name,
				IPAddress: clientResponse.IPAddress,
			}
			channelClient <- client
		}
	}()
	loginResponseWithClient := models.LoginResponseWithClient{
		LoginResponse: *loginResponse,
		Client:        <-channelClient,
	}

	ctx.JSON(http.StatusOK, loginResponseWithClient)
}
func (h *authHandler) Register(ctx *gin.Context) {
	var req models.RegisterUserRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	registerUserResponse, err := h.service.RegisterUser(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, registerUserResponse)
}

func (h *authHandler) RefreshToken(ctx *gin.Context) {
	var req models.RefreshTokenRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshTokenResponse, err := h.service.RefreshToken(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, refreshTokenResponse)
}
