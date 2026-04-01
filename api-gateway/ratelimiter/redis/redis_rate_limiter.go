package redis

import (
	"context"
	"errors"
	"gaman-microservice/api-gateway/ratelimiter"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	rateLimitTTL    = 10 * time.Minute
	rateLimitPrefix = "apigw:ratelimit:"
)

func genKey(id, method string) string {
	return rateLimitPrefix + id + ":" + method
}

type redisRateLimiter struct {
	redis *redis.Client
}

func NewRedisRateLimiter(redis *redis.Client) ratelimiter.RateLimiter {
	return &redisRateLimiter{
		redis: redis,
	}
}

func (r redisRateLimiter) ValidateRateLimit(ctx context.Context, limit int32, id string, method string) (*ratelimiter.RateLimitResult, error) {
	key := genKey(id, method)
	l := log.Ctx(ctx).With().
		Str("module", "redisRateLimiter").
		Str("key", key).
		Logger()

	// set with default value
	out := ratelimiter.RateLimitResult{
		Allow:          true,
		QuotaLeft:      limit,
		QuotaResetTime: time.Now().Add(rateLimitTTL),
	}

	currentStr, err := r.redis.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		l.Warn().Err(err).Msg("error getting existing value")
		return nil, err
	}
	current, err := strconv.Atoi(currentStr)
	if err != nil {
		current = 0
		l.Warn().
			Err(err).
			Str("value", currentStr).
			Msg("error parsing existing value fallback to 0")
	}
	out.QuotaLeft = limit - int32(current)
	out.Allow = out.QuotaLeft > 0
	if out.Allow {
		out.QuotaResetTime = time.Now().Add(rateLimitTTL)
		// increment and re-set the ttl
		_, err := r.redis.TxPipelined(ctx, func(tx redis.Pipeliner) error {
			tx.Incr(ctx, key)
			tx.Expire(ctx, key, rateLimitTTL)
			return nil
		})
		if err != nil {
			l.Err(err).Any("out", out).Msg("error incrementing and setting ttl")
			return nil, err
		}
		return &out, nil
	}

	// not allowed get existing TTL
	ttl, err := r.redis.TTL(ctx, key).Result()
	if err != nil {
		l.Err(err).Msg("error getting existing TTL fallback to default")
		ttl = rateLimitTTL
	}
	out.QuotaResetTime = time.Now().Add(ttl)
	return &out, nil
}
