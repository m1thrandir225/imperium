package server

import "context"

type StatusHandler interface {
	GetStatus(ctx context.Context)
}

type statusHandler struct{}

func NewStatusHandler() StatusHandler {
	return &statusHandler{}
}

func (h *statusHandler) GetStatus(ctx context.Context) {

}
