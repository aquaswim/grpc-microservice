package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"gaman-microservice/api-gateway/constant"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func RequestIdMiddleware() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestId()
			}

			// Add request ID to response header
			w.Header().Set("X-Request-ID", requestID)

			// Add request ID to context
			ctx := r.Context()
			ctx = context.WithValue(ctx, constant.CtxKeyRequestID, requestID)
			r = r.WithContext(ctx)

			next(w, r, pathParams)
		}
	}
}

func generateRequestId() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "--getrandom-error--"
	}
	return hex.EncodeToString(b)
}
