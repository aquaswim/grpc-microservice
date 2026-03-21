package main

import (
	"context"
	"gaman-microservice/api-gateway/config"
	userv1 "gaman-microservice/api-gateway/gen/user/v1"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "Authorization":
				return key, true
			default:
				return runtime.DefaultHeaderMatcher(key)
			}
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = userv1.RegisterUserServiceHandlerFromEndpoint(ctx, mux, cfg.UserSvcAddr, opts)
	if err != nil {
		log.Panicf("failed to register gateway: %v", err)
	}

	// start server
	log.Printf("Starting HTTP server on address %s", cfg.ListenAddr)
	err = http.ListenAndServe(cfg.ListenAddr, mux)
	if err != nil {
		log.Printf("http.ListenAndServe error %+v", err)
	}
}
