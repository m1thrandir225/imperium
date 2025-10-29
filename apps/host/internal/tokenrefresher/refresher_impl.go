package tokenrefresher

import (
	"context"
	"log"
	"time"

	"github.com/m1thrandir225/imperium/apps/host/internal/httpclient"
)

type AuthTokenRefresher struct {
	getter    httpclient.TokenGetter
	refresher httpclient.TokenRefresher
	stop      chan struct{}
}

func NewTokenRefresher(getter httpclient.TokenGetter, refresher httpclient.TokenRefresher) (Refresher, error) {
	if getter == nil {
		return nil, ErrInvalidTokenGetter
	}
	if refresher == nil {
		return nil, ErrInvalidTokenRefrehser
	}
	return &AuthTokenRefresher{
		getter:    getter,
		refresher: refresher,
		stop:      make(chan struct{}),
	}, nil
}

func (r *AuthTokenRefresher) Start(ctx context.Context) {
	t := time.NewTicker(1 * time.Minute)
	go func() {
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return

			case <-r.stop:
				return
			case <-t.C:
				expiresAt := r.getter.GetAccessTokenExpiresAt()
				if time.Until(expiresAt) < 5*time.Minute {
					c, cancel := context.WithTimeout(ctx, 10*time.Second)
					if err := r.refresher.RefreshToken(c); err != nil {
						log.Printf("failed to refresh token: %v", err)
					}
					cancel()
				}
			}
		}

	}()
}

func (r *AuthTokenRefresher) Stop() {
	select {
	case <-r.stop:
	default:
		close(r.stop)
	}
}
