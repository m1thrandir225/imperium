package tokenrefresher

import "context"

type Refresher interface {
	Start(ctx context.Context)
	Stop()
}
