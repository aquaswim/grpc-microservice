package ratelimiter

import (
	"context"
	"time"
)

type RateLimitResult struct {
	Allow          bool
	QuotaLeft      int32
	QuotaResetTime time.Time
}

type RateLimiter interface {
	ValidateRateLimit(ctx context.Context, limit int32, id string, method string) (*RateLimitResult, error)
}
