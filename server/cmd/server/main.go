package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/serbanmarti/go-grpc/env"
	"github.com/serbanmarti/go-grpc/proto_gen/crud/v1/crudv1connect"
	"github.com/serbanmarti/go-grpc/proto_gen/stream/v1/streamv1connect"
	"github.com/serbanmarti/go-grpc/server/interceptor"
	"github.com/serbanmarti/go-grpc/server/service"
)

func main() {
	// Get the environment configuration
	environment := env.GetEnvironment()

	// Initialize the logger
	var logger *zap.Logger
	var err error
	switch environment.Environment {
	case "production":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatalf("Failed to create prod logger: %v\n", err)
		}
	case "development":
		fallthrough
	default:
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatalf("Failed to create dev logger: %v\n", err)
		}
	}
	zap.ReplaceGlobals(logger)

	// Instantiate the interceptors
	interceptors := connect.WithInterceptors(
		interceptor.NewLoggerInterceptor(),
		interceptor.NewAuthInterceptor(),
		interceptor.NewRecoveryInterceptor(),
	)

	// Create the server mux
	mux := http.NewServeMux()

	// Register the proto services
	mux.Handle(crudv1connect.NewCrudServiceHandler(&service.CrudService{
		Data:  make(map[string]string),
		Mutex: sync.RWMutex{},
	}, interceptors))
	mux.Handle(streamv1connect.NewStreamServiceHandler(&service.StreamService{}, interceptors))

	// Register the reflection service on the server
	reflector := grpcreflect.NewStaticReflector(
		crudv1connect.CrudServiceName,
		streamv1connect.StreamServiceName,
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Create the server
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", environment.Port),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Create a channel to listen for OS signals
	openConnsClosed := make(chan struct{})
	go func() {
		// Wait for a signal
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM) // Catch SIGINT (Ctrl+C) and SIGTERM
		<-sigint

		// Received a signal, gracefully shut down
		zap.L().Info("Shutting down server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout
			zap.L().Error(fmt.Sprintf("HTTP server shutdown error: %v", err))
		}

		// Signal that open connections are closed
		close(openConnsClosed)
	}()

	// Start the server
	zap.L().Info(fmt.Sprintf("Starting server and listening at %s...", srv.Addr))
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener
		zap.L().Error(fmt.Sprintf("HTTP server listen error: %v", err))
	} else {
		zap.L().Info("Server shut down!")
	}

	// Wait for ongoing connections to close (from the shutdown signal received)
	<-openConnsClosed
}
