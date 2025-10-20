package host

import "context"

type StatusManager interface {
	Start(ctx context.Context)
	Stop()
	UpdateStatus(status Status)
	SendStatusUpdate(ctx context.Context, status Status)
}
