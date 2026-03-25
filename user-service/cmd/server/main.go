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

type AppDep struct {
	Cfg               *config.Config              `container:"type"`
	UserAuthHandler   *usergrpc.UserAuthHandler   `container:"type"`
	UserManageHandler *usergrpc.UserManageHandler `container:"type"`
}

func main() {
	// Initialize Container
	c := container.Init()

	appDep := AppDep{}
	err := c.Fill(&appDep)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to fill app dependencies")
	}

	// Create TCP listener
	lis, err := net.Listen("tcp", appDep.Cfg.TcpListenerUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	// Create gRPC server
	s := usergrpc.NewServer()

	// Register Services
	userv1.RegisterManageServiceServer(s, appDep.UserManageHandler)
	userv1.RegisterAuthServiceServer(s, appDep.UserAuthHandler)

	// Register reflection service on gRPC server
	reflection.Register(s)

	log.Info().Str("url", appDep.Cfg.TcpListenerUrl).Msg("User Service is running")
	if err := s.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
