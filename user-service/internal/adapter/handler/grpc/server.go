package grpc

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const metadataKeyRequestId = "x-request-id"

type requestIdKey struct{}

func UnaryRequestIdInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if requestIds := md.Get(metadataKeyRequestId); len(requestIds) > 0 {
			ctx = context.WithValue(ctx, requestIdKey{}, requestIds[0])
		}
	}

	// register logger
	ctx = registerLogger(ctx, info.FullMethod)

	return handler(ctx, req)
}

func StreamRequestIdInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if requestIds := md.Get(metadataKeyRequestId); len(requestIds) > 0 {
			ctx = context.WithValue(ctx, requestIdKey{}, requestIds[0])
		}
	}

	// register logger
	ctx = registerLogger(ctx, info.FullMethod)

	wrapped := &wrappedServerStream{ServerStream: ss, ctx: ctx}
	return handler(srv, wrapped)
}

type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

func UnaryLoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("request started")

	resp, err := handler(ctx, req)

	duration := time.Since(start)
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, err.Error())
	}

	logger.Info().
		Dur("duration", duration).
		Str("status", st.Code().String()).
		Msg("request finished")

	return resp, err
}

func StreamLoggingInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	ctx := ss.Context()
	logger := zerolog.Ctx(ctx)

	logger.Info().Msg("stream request started")

	err := handler(srv, ss)

	duration := time.Since(start)
	st, ok := status.FromError(err)
	if !ok {
		st = status.New(codes.Unknown, err.Error())
	}

	logger.Info().
		Dur("duration", duration).
		Str("status", st.Code().String()).
		Msg("stream request finished")

	return err
}

func UnaryRecoveryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			logger := zerolog.Ctx(ctx)
			logger.Error().
				Interface("panic", r).
				Str("stack", string(debug.Stack())).
				Msg("panic recovered in unary interceptor")
			err = status.Errorf(codes.Internal, "internal server error: %v", r)
		}
	}()
	return handler(ctx, req)
}

func StreamRecoveryInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	defer func() {
		if r := recover(); r != nil {
			logger := zerolog.Ctx(ss.Context())
			logger.Error().
				Interface("panic", r).
				Str("stack", string(debug.Stack())).
				Msg("panic recovered in stream interceptor")
			err = status.Errorf(codes.Internal, "internal server error: %v", r)
		}
	}()
	return handler(srv, ss)
}

func NewServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryRequestIdInterceptor,
			UnaryRecoveryInterceptor,
			UnaryLoggingInterceptor,
		),
		grpc.ChainStreamInterceptor(
			StreamRequestIdInterceptor,
			StreamRecoveryInterceptor,
			StreamLoggingInterceptor,
		),
	)
}

func getRequestIdFromContext(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIdKey{}).(string); ok {
		return requestId
	}
	return "--request-id-err--"
}

func registerLogger(ctx context.Context, methodName string) context.Context {
	l := log.With().
		Str("reqId", getRequestIdFromContext(ctx)).
		Str("method", methodName).
		Logger()

	return l.WithContext(ctx)
}
