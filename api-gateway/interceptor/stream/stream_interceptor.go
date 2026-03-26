package stream

import (
	"context"
	"gaman-microservice/api-gateway/constant"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GatewayInterceptor() []grpc.StreamClientInterceptor {
	return []grpc.StreamClientInterceptor{
		RequestIdInterceptor(),
		LogInterceptor(),
	}
}

func LogInterceptor() grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		l := zerolog.Ctx(ctx)
		l.Debug().Msgf("[stream forward] %s", method)

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func RequestIdInterceptor() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if requestID, ok := ctx.Value(constant.CtxKeyRequestID).(string); ok {
			ctx = metadata.AppendToOutgoingContext(ctx, constant.MetadataKeyRequestID, requestID)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}
