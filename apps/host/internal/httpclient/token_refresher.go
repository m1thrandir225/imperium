package httpclient

import "context"

type TokenRefresher interface {
	RefreshToken(ctx context.Context) error
}
