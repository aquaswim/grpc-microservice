package unary

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/interceptor/utils"
	"gaman-microservice/api-gateway/methodmetamap"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GatewayInterceptor(methodMap methodmetamap.MethodMetaMap, authClient userv1.AuthServiceClient) []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{
		RequestIdInterceptor(),
		AuthInterceptor(methodMap, authClient),
		LogInterceptor(),
	}
}

func LogInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		l := zerolog.Ctx(ctx)

		l.Debug().Msgf("[unary forward] %s", method)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func RequestIdInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if requestID, ok := ctx.Value(constant.CtxKeyRequestID).(string); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, constant.MetadataKeyRequestID, requestID)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func AuthInterceptor(methodMap methodmetamap.MethodMetaMap, authClient userv1.AuthServiceClient) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		l := log.Ctx(ctx)
		if methodMeta, ok := methodMap.Get(method); ok {
			if methodMeta.NeedAuth {
				// check for auth
				tokenData, err := utils.ValidateTokenFromContext(ctx, authClient)
				if err != nil {
					return err
				}
				l.UpdateContext(func(c zerolog.Context) zerolog.Context {
					return c.Str("user_id", tokenData.Data.Id)
				})

				ctx = metadata.AppendToOutgoingContext(ctx,
					constant.MetadataKeyUserId, tokenData.Data.GetId(),
					constant.MetadataKeyUsername, tokenData.Data.GetUsername(),
				)
			}
		} else {
			// method map not found
			l.Warn().
				Str("method", method).
				Msg("method not found in methodMetaMap, gateway maybe not configured properly")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
