package main

import (
	"context"
	"gaman-microservice/api-gateway/config"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/interceptor/stream"
	"gaman-microservice/api-gateway/interceptor/unary"
	"gaman-microservice/api-gateway/middleware"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if cfg.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		log.Warn().Msgf("Pretty logging is enabled, this must only be used in local!")
	}

	mux := runtime.NewServeMux(
		runtime.WithMiddlewares(middleware.GatewayMiddleware()...),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(unary.GatewayInterceptor()...),
		grpc.WithChainStreamInterceptor(stream.GatewayInterceptor()...),
	}
	err = userv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.UserSvcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register gateway")
	}

	// start server
	log.
		Info().
		Str("address", cfg.ListenAddr).
		Msg("Starting HTTP server address")
	err = http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		log.Info().Err(err).Msg("http.ListenAndServe error")
	}
}
