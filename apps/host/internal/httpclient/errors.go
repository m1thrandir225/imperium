package httpclient

import "errors"

var (
	InvalidBaseUrl        = errors.New("invalid base URL")
	InvalidTokenGetter    = errors.New("invalid TokenGetter")
	InvalidTokenRefresher = errors.New("invalid TokenRefresher")
)
