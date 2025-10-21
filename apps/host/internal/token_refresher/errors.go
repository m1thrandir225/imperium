package tokenrefresher

import "errors"

var (
	ErrInvalidTokenGetter    = errors.New("invalid TokenGetter")
	ErrInvalidTokenRefrehser = errors.New("invalid TokenRefresher")
)
