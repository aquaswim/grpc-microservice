package ratelimiter

import (
	"context"
	"time"
)

type RateLimitResult struct {
	Allow          bool
	QuotaLeft      int32
	QuotaResetLeft time.Duration
}

func (r RateLimitResult) QuotaResetTime() time.Time {
	return time.Now().Add(r.QuotaResetLeft)
}

type RateLimiter interface {
	ValidateRateLimit(ctx context.Context, limit int32, id string, method string) (*RateLimitResult, error)
}
