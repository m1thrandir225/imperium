package internal

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
)

type HTTPHandler struct {
	config        *config.Config
	authService   *AuthService
	hostService   *HostService
	clientService *ClientService
}

func NewHTTPHandler(
	config *config.Config,
	authService *AuthService,
	hostService *HostService,
	clientService *ClientService,
) *HTTPHandler {
	return &HTTPHandler{
		config:        config,
		authService:   authService,
		hostService:   hostService,
		clientService: clientService,
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

	channelClient := make(chan ClientDTO)
	// upsert client
	go func() {
		hostname, err := GetHostname()
		if err != nil {
			hostname = "unknown-hostname"
		}
		ip, err := GetIPV4Address()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		clientResponse, err := h.clientService.RegisterOrUpdateClient(ctx, RegisterClientRequest{
			Name:      hostname,
			IPAddress: ip,
		}, loginResponse.AccessToken)
		if err != nil {
			log.Println("error upserting client", err)
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		} else {
			log.Println("client upserted successfully")
			client := ClientDTO{
				ID:        clientResponse.ID,
				Name:      clientResponse.Name,
				IPAddress: clientResponse.IPAddress,
			}
			channelClient <- client
		}
	}()
	loginResponseWithClient := LoginResponseWithClient{
		LoginResponse: *loginResponse,
		Client:        <-channelClient,
	}

	ctx.JSON(http.StatusOK, loginResponseWithClient)
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

func (h *HTTPHandler) GetHost(ctx *gin.Context) {
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

	host, err := h.hostService.GetHost(ctx, uriID.HostID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, host)
}

func (h *HTTPHandler) GetConfigStatus(ctx *gin.Context) {
	configured := config.IsConfigured()

	ctx.JSON(http.StatusOK, gin.H{
		"configured": configured,
	})
}

func (h *HTTPHandler) SetupConfig(ctx *gin.Context) {
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

	ctx.JSON(http.StatusOK, messageResponse("Config setup successfully"))
}

func (h *HTTPHandler) GetHostPrograms(ctx *gin.Context) {

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

	programs, err := h.hostService.GetHostPrograms(ctx, uriID.HostID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, programs)
}
