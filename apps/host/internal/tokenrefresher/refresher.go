// Package tokenrefresher provides a definition and implementation of the a
// time-based JWT/PASETO token refresher that automatically refreshes a given
// access_token
package tokenrefresher

import "context"

type Refresher interface {
	Start(ctx context.Context)
	Stop()
}
