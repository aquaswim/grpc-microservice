package middleware

import (
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
)

func RecoverMiddleware(next runtime.HandlerFunc) runtime.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		defer func() {
			if err := recover(); err != nil {
				zerolog.
					Ctx(r.Context()).
					Error().
					Stack().
					Any("Error", err).
					Msg("panic recovery")
				w.WriteHeader(http.StatusBadGateway)
				_, _ = w.Write([]byte(http.StatusText(http.StatusBadGateway)))
			}
		}()
		next(w, r, pathParams)
	}
}
