package main

import (
	userv1 "gaman-microservice/user-service/gen/user/v1"
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/container"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize Container
	c := container.Init()

	var cfg *config.Config
	if err := c.Resolve(&cfg); err != nil {
		log.Fatal().Err(err).Msg("failed to resolve config")
	}

	var userHandler *usergrpc.UserHandler
	if err := c.Resolve(&userHandler); err != nil {
		log.Fatal().Err(err).Msg("failed to resolve user handler")
	}

	// Create TCP listener
	lis, err := net.Listen("tcp", cfg.TcpListenerUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	// Create gRPC server
	s := usergrpc.NewServer()

	// Register UserService
	userv1.RegisterUserServiceServer(s, userHandler)

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Info().Str("url", cfg.TcpListenerUrl).Msg("User Service is running")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
