package httpclient

import "time"

type TokenGetter interface {
	GetAccessToken() string
	IsAccessTokenExpired() bool
	GetAccessTokenExpiresAt() time.Time
}
