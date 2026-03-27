package grpcInterceptorUtil

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/methodmetamap"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/metadata"
)

type MethodMetaProcessor struct {
	mm         methodmetamap.MethodMetaMap
	authClient userv1.AuthServiceClient
}

func NewMethodMetaProcessor(
	methodMetaMap methodmetamap.MethodMetaMap,
	authClient userv1.AuthServiceClient,
) *MethodMetaProcessor {
	return &MethodMetaProcessor{
		mm:         methodMetaMap,
		authClient: authClient,
	}
}

func (mmp *MethodMetaProcessor) ProcessMethodMeta(ctx context.Context, method string) (context.Context, error) {
	l := zerolog.Ctx(ctx)

	meta, ok := mmp.mm.Get(method)
	if !ok {
		// method map not found
		l.Warn().
			Str("method", method).
			Msg("method not found in methodMetaMap, gateway maybe not configured properly")
		return ctx, nil
	}
	if meta.NeedAuth {
		// check auth
		tokenData, err := ValidateTokenFromContext(ctx, mmp.authClient)
		if err != nil {
			return nil, err
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

	return ctx, nil
}
