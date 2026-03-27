package main

import (
	"context"
	"gaman-microservice/api-gateway/config"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"gaman-microservice/api-gateway/interceptor/stream"
	"gaman-microservice/api-gateway/interceptor/unary"
	grpcInterceptorUtil "gaman-microservice/api-gateway/interceptor/utils"
	"gaman-microservice/api-gateway/methodmetamap"
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

	// method meta map
	methodMap, err := methodmetamap.GetMethodMetaFromFileDesc(
		userv1.File_user_v1_auth_proto,
		userv1.File_user_v1_manage_proto,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load method meta map")
	}

	// auth client
	authClientConn, err := grpc.NewClient(
		cfg.UserSvcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create auth client")
	}
	defer func() {
		err := authClientConn.Close()
		if err != nil {
			log.Err(err).Msg("failed to close auth client connection")
		}
	}()
	authSvcClient := userv1.NewAuthServiceClient(authClientConn)

	if authSvcClient == nil {
		log.Fatal().Msg("failed to create auth service client")
	}

	mmProcessor := grpcInterceptorUtil.NewMethodMetaProcessor(methodMap, authSvcClient)

	mux := runtime.NewServeMux(
		runtime.WithMiddlewares(middleware.GatewayMiddleware()...),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(unary.GatewayInterceptor(mmProcessor)...),
		grpc.WithChainStreamInterceptor(stream.GatewayInterceptor(mmProcessor)...),
	}
	err = userv1.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, cfg.UserSvcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register service user/auth to gateway")
	}
	err = userv1.RegisterManageServiceHandlerFromEndpoint(ctx, mux, cfg.UserSvcAddr, opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to register service user/manage to gateway")
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
