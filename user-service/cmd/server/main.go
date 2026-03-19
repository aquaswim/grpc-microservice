package main

import (
	"fmt"
	"log"
	"net"

	userv1 "gaman-microservice/user-service/gen/user/v1"
	usergrpc "gaman-microservice/user-service/internal/adapter/handler/grpc"
	"gaman-microservice/user-service/internal/infrastructure/config"
	"gaman-microservice/user-service/internal/infrastructure/container"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize Container
	c := container.Init()

	var cfg *config.Config
	if err := c.Resolve(&cfg); err != nil {
		log.Fatalf("failed to resolve config: %v", err)
	}

	var userHandler *usergrpc.UserHandler
	if err := c.Resolve(&userHandler); err != nil {
		log.Fatalf("failed to resolve user handler: %v", err)
	}

	// Create TCP listener
	lis, err := net.Listen("tcp", cfg.TcpListenerUrl)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()

	// Register UserService
	userv1.RegisterUserServiceServer(s, userHandler)

	// Register reflection service on gRPC server
	reflection.Register(s)

	fmt.Printf("User Service is running on %s...\n", cfg.TcpListenerUrl)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
