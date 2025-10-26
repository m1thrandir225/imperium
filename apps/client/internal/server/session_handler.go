package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m1thrandir225/imperium/apps/client/config"
	"github.com/m1thrandir225/imperium/apps/client/internal/models"
	"github.com/m1thrandir225/imperium/apps/client/internal/services"
)

type SessionHandler interface {
	CreateSession(ctx *gin.Context)
	GetSession(ctx *gin.Context)
	StartSession(ctx *gin.Context)
	EndSession(ctx *gin.Context)
	CancelSession(ctx *gin.Context)
}

type sessionHandler struct {
	config  *config.Config
	service services.SessionService
}

func NewSessionHandler(
	cfg *config.Config,
	service services.SessionService,
) SessionHandler {
	return &sessionHandler{
		config:  cfg,
		service: service,
	}
}

func (h *sessionHandler) CreateSession(ctx *gin.Context) {
	token := GetAuthToken(ctx)
	if token == "" {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}

	var req models.CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.service.CreateSession(ctx, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *sessionHandler) GetSession(ctx *gin.Context) {
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

	sessionResponse, err := h.service.GetSession(ctx, sessionID, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *sessionHandler) StartSession(ctx *gin.Context) {
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

	var req models.StartSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.service.StartSession(ctx, sessionID, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *sessionHandler) EndSession(ctx *gin.Context) {
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

	var req models.EndSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	sessionResponse, err := h.service.EndSession(ctx, sessionID, req, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}

func (h *sessionHandler) CancelSession(ctx *gin.Context) {
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

	sessionResponse, err := h.service.CancelSession(ctx, sessionID, reason, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, sessionResponse)
}
