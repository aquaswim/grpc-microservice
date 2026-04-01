package middleware

import (
	"context"
	"gaman-microservice/api-gateway/constant"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
)

func AppendIpMiddleware() runtime.Middleware {
	return func(next runtime.HandlerFunc) runtime.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Ctx(r.Context()).Err(err).Caller().Msg("Failed to parse IP from request")
				ip = "unknown"
			}

			ctx := context.WithValue(r.Context(), constant.CtxKeyIP, ip)

			next(w, r.WithContext(ctx), pathParams)
		}
	}
}
