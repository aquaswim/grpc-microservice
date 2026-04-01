package grpcInterceptorUtil

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/methodmetamap"
	"gaman-microservice/api-gateway/ratelimiter"

	"github.com/rs/zerolog"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
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

		ctxIp, ok := ctx.Value(constant.CtxKeyIP).(string)
		if !ok {
			ctxIp = "unknown"
		}

		rateLimitResult, err := mmp.rateLimiter.ValidateRateLimit(ctx, meta.RateLimit, ctxIp, method)
		if err != nil {
			return ctx, err
		}
		if !rateLimitResult.Allow {
			l.Error().
				Int32("rate_limit", meta.RateLimit).
				Int32("quota_left", rateLimitResult.QuotaLeft).
				Time("quota_reset_at", rateLimitResult.QuotaResetTime()).
				Msg("rate limit exceeded")

			st, errx := status.New(codes.ResourceExhausted, "rate limit exceeded").
				WithDetails(&errdetails.RetryInfo{
					RetryDelay: durationpb.New(rateLimitResult.QuotaResetLeft),
				})
			if errx != nil {
				l.Error().Err(errx).Msg("failed to add rate limit details")
			}

			return ctx, st.Err()
		}
	}

	return ctx, nil
}
