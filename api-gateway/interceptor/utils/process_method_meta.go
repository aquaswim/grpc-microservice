package grpcInterceptorUtil

import (
	"context"
	"fmt"
	"gaman-microservice/api-gateway/constant"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/methodmetamap"
	"gaman-microservice/api-gateway/ratelimiter"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

type MethodMetaProcessor struct {
	mm          methodmetamap.MethodMetaMap
	authClient  userv1.AuthServiceClient
	rateLimiter ratelimiter.RateLimiter
}

func NewMethodMetaProcessor(
	methodMetaMap methodmetamap.MethodMetaMap,
	authClient userv1.AuthServiceClient,
	rl ratelimiter.RateLimiter,
) *MethodMetaProcessor {
	return &MethodMetaProcessor{
		mm:          methodMetaMap,
		authClient:  authClient,
		rateLimiter: rl,
	}
}

func (mmp *MethodMetaProcessor) ProcessMethodMeta(ctx context.Context, method string) (context.Context, error) {
	l := zerolog.Ctx(ctx).With().Str("method", method).Logger()

	meta, ok := mmp.mm.Get(method)
	if !ok {
		// method map not found
		l.Warn().Msg("method not found in methodMetaMap, gateway maybe not configured properly")
		return ctx, nil
	}
	if meta.NeedAuth {
		// check auth
		tokenData, err := ValidateTokenFromContext(ctx, mmp.authClient)
		if err != nil {
			return ctx, err
		}
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str(constant.MetadataKeyUserId, tokenData.Data.GetId()).
				Str(constant.MetadataKeyUsername, tokenData.Data.GetUsername())
		})

		ctx = metadata.AppendToOutgoingContext(ctx,
			constant.MetadataKeyUserId, tokenData.Data.GetId(),
			constant.MetadataKeyUsername, tokenData.Data.GetUsername(),
		)
	}
	if meta.RateLimit > 0 {
		l.Debug().Int32("rate_limit", meta.RateLimit).Msg("check rate limit")
		// todo userID from IP
		rateLimitResult, err := mmp.rateLimiter.ValidateRateLimit(ctx, meta.RateLimit, "x", method)
		if err != nil {
			return ctx, err
		}
		// todo return rate limit info as header
		if !rateLimitResult.Allow {
			l.Error().
				Int32("rate_limit", meta.RateLimit).
				Int32("quota_left", rateLimitResult.QuotaLeft).
				Time("quota_reset_at", rateLimitResult.QuotaResetTime).
				Msg("rate limit exceeded")
			return ctx, &runtime.HTTPStatusError{
				HTTPStatus: http.StatusTooManyRequests,
				Err:        fmt.Errorf("rate limit exceeded"),
			}
		}
	}

	return ctx, nil
}
