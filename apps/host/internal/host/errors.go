package host

import "errors"

var (
	ErrInvalidHostID            = errors.New("invalid hostID")
	ErrInvalidAuthServerBaseURL = errors.New("invalid authServerBaseURL")
	ErrInvalidHttpClient        = errors.New("invalid httpclient.Client")
	ErrInvalidSessionService    = errors.New("invalid sessionService")
)
