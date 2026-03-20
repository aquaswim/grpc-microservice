package main

import (
	"context"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	userSvcEndpoint = "localhost:50051"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userSvcEndpoint, opts)
	if err != nil {
		log.Panicf("failed to register gateway: %v", err)
	}

	// start server
	log.Printf("Starting HTTP server on port %s", ":8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Printf("http.ListenAndServe error %+v", err)
	}
}
