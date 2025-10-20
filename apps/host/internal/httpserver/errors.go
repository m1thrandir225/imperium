package httpserver

import "errors"

var (
	InvalidSessionService = errors.New("invalid SessionService")
	InvalidEventBus       = errors.New("invalid EventBus")
)
