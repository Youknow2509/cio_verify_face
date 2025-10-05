package ratelimit

import (
	"context"

	"github.com/youknow2509/cio_verify_face/server/service_ws_delivery/internal/domain/model"
)

// ===============================================
// Limiter interface
// ===============================================
type ILimiter interface {
	Check(ctx context.Context, key string) (int, bool, error)
	Create(ctx context.Context, key string) error
	Upgrade(ctx context.Context, key string, val int) error
	Allow(ctx context.Context, key string) (model.Verdict, error)
}
