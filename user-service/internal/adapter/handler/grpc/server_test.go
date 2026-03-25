package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Define a handler that panics
	unaryHandler := func(ctx context.Context, req any) (any, error) {
		panic("test panic")
	}

	// Create a dummy UnaryServerInfo
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.TestService/TestPanic",
	}

	// This should NOT panic now because we are using the recovery interceptor
	resp, err := UnaryRecoveryInterceptor(context.Background(), nil, info, unaryHandler)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected gRPC status error, got %T", err)
	}

	if st.Code() != codes.Internal {
		t.Errorf("expected codes.Internal, got %v", st.Code())
	}

	if resp != nil {
		t.Errorf("expected nil response, got %v", resp)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	unaryHandler := func(ctx context.Context, req any) (any, error) {
		return "ok", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.TestService/TestLogging",
	}

	resp, err := UnaryLoggingInterceptor(context.Background(), nil, info, unaryHandler)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp != "ok" {
		t.Errorf("expected 'ok', got %v", resp)
	}
}
