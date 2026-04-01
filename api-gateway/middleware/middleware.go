package middleware

import (
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func GatewayMiddleware() []runtime.Middleware {
	return []runtime.Middleware{
		RequestIdMiddleware(),
		LoggerMiddleware(),
		AppendIpMiddleware(),
		RecoverMiddleware,
	}
}
