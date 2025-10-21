package auth

import "errors"

var (
	ErrInvalidAuthServiceBaseURL = errors.New("invalid authServiceBaseURL")
	ErrInvalidHttpClient         = errors.New("invalid httpclient.Client")
)
