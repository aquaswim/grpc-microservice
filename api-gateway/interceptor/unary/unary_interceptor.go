package unary

import (
	"context"
	"gaman-microservice/api-gateway/constant"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GatewayInterceptor() []grpc.UnaryClientInterceptor {
	return []grpc.UnaryClientInterceptor{
		RequestIdInterceptor(),
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
