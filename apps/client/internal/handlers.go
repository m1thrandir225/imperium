package internal

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
)

type HTTPHandler struct {
	config         *config.Config
	authService    *AuthService
	hostService    *HostService
	clientService  *ClientService
	sessionService *SessionService
}

func NewHTTPHandler(
	config *config.Config,
	authService *AuthService,
	hostService *HostService,
	clientService *ClientService,
	sessionService *SessionService,
) *HTTPHandler {
	return &HTTPHandler{
		config:         config,
		authService:    authService,
		hostService:    hostService,
		clientService:  clientService,
		sessionService: sessionService,
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

	h.updateServicesWithNewURL(req.AuthServerBaseURL)

	ctx.JSON(http.StatusOK, messageResponse("Config setup successfully"))
}

// updateServicesWithNewURL updates all services with the new auth server base URL
func (h *HTTPHandler) updateServicesWithNewURL(authServerBaseURL string) {
	apiURL := fmt.Sprintf("%s/api/v1", authServerBaseURL)

	h.authService.UpdateAuthServerBaseURL(fmt.Sprintf("%s/auth", apiURL))
	h.hostService.UpdateHostServerBaseURL(fmt.Sprintf("%s/hosts", apiURL))
	h.clientService.UpdateBaseURL(fmt.Sprintf("%s/clients", apiURL))
	h.sessionService.UpdateBaseURL(fmt.Sprintf("%s/sessions", apiURL))
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

func (h *HTTPHandler) GetClientInfo(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	hostname, ip, err := h.clientService.GetClientInfo()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	clientReq := RegisterClientRequest{
		Name:      hostname,
		IPAddress: ip,
	}

	clientResponse, err := h.clientService.RegisterOrUpdateClient(ctx, clientReq, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, clientResponse)
}

func (h *HTTPHandler) CreateSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	var req CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.sessionService.CreateSession(ctx, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *HTTPHandler) GetSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("session id is required")))
		return
	}

	sessionResponse, err := h.sessionService.GetSession(ctx, sessionID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *HTTPHandler) StartSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("session id is required")))
		return
	}

	var req StartSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.sessionService.StartSession(ctx, sessionID, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *HTTPHandler) EndSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("session id is required")))
		return
	}

	var req EndSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.sessionService.EndSession(ctx, sessionID, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *HTTPHandler) CancelSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	sessionID := ctx.Param("sessionId")
	if sessionID == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("session id is required")))
		return
	}

	reason := ctx.Query("reason")
	if reason == "" {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("reason is required")))
		return
	}

	sessionResponse, err := h.sessionService.CancelSession(ctx, sessionID, reason, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}
