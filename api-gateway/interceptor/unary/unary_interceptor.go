package unary

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	grpcInterceptorUtil "gaman-microservice/api-gateway/interceptor/utils"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GatewayInterceptor(methodMetaProcessor *grpcInterceptorUtil.MethodMetaProcessor) []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{
		RequestIdInterceptor(),
		MethodMetaProcessorInterceptor(methodMetaProcessor),
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

func MethodMetaProcessorInterceptor(processor *grpcInterceptorUtil.MethodMetaProcessor) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		newCtx, err := processor.ProcessMethodMeta(ctx, method)
		if err != nil {
			return err
		}
		return invoker(newCtx, method, req, reply, cc, opts...)
	}
}
