package tokenrefresher

import "errors"

var (
	ErrInvalidTokenGetter    = errors.New("invalid TokenGetter")
	ErrInvalidTokenRefresher = errors.New("invalid TokenRefresher")
)
