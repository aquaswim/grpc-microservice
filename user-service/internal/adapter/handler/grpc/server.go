package grpc

import (
	"context"
	"errors"
	commonv1 "gaman-microservice/user-service/gen/common/v1"
	appError "gaman-microservice/user-service/internal/domain/app_error"
	"runtime/debug"
	"time"

	"github.com/joomcode/errorx"
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
		Err(err).
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
		Err(err).
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

func getRequestIdFromContext(ctx context.Context) string {
	if requestId, ok := ctx.Value(requestIdKey{}).(string); ok {
		return requestId
	}
	return "--request-id-err--"
}

func registerLogger(ctx context.Context, methodName string) context.Context {
	userCtx, err := getAuthDataFromCtx(ctx)
	if err != nil {
		// user not logged in
		userCtx = &commonv1.TokenPayload{
			Id: "guest",
		}
	}
	l := log.With().
		Str("reqId", getRequestIdFromContext(ctx)).
		Str("method", methodName).
		Str("user_id", userCtx.GetId()).
		Logger()

	return l.WithContext(ctx)
}

func toGrpcError(code codes.Code, err error) error {
	var appErr *errorx.Error
	if !errors.As(err, &appErr) {
		return status.Errorf(code, "%s", err)
	}
	return status.Errorf(code, "%s", appErr.Message())
}

func errorMapper(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}

	zerolog.Ctx(ctx).Error().Err(err).Msgf("Error: %+v", err)

	errType := appError.Switch(err)

	switch errType {
	case appError.ErrNotFound:
		return toGrpcError(codes.NotFound, err)
	case appError.ErrValidation:
		return toGrpcError(codes.InvalidArgument, err)
	case appError.ErrUnauthorized:
		return toGrpcError(codes.Unauthenticated, err)
	default:
		return toGrpcError(codes.Internal, err)
	}
}

func StreamErrorMappingInterceptor(srv any, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	return errorMapper(ss.Context(), err)
}

func UnaryErrorMappingInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		return nil, errorMapper(ctx, err)
	}
	return resp, nil
}

func NewServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryRequestIdInterceptor,
			UnaryRecoveryInterceptor,
			UnaryLoggingInterceptor,
			UnaryErrorMappingInterceptor,
		),
		grpc.ChainStreamInterceptor(
			StreamRequestIdInterceptor,
			StreamRecoveryInterceptor,
			StreamLoggingInterceptor,
			StreamErrorMappingInterceptor,
		),
	)
}
