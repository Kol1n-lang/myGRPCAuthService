package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"authService/github.com/authService/api"
	"authService/internal/config"
	httpServe "authService/internal/infrastructure/http"
	"authService/internal/infrastructure/implementations/broker"
	"authService/internal/infrastructure/implementations/postgres"
	"authService/internal/infrastructure/middleware"
	"authService/internal/monitoring"
	"authService/internal/service"
	"google.golang.org/grpc"
)

func main() {
	monitoring.InitMetrics()
	cfg := config.Init()
	log.Println("Config initialized")

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	err := config.RunMigrations(cfg.DB.DBUrl())
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.MetricsInterceptor("auth-service"),
		),
	)

	db, err := config.CreateDBConnection(cfg.DB.DBUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	userRepository := postgres.NewUserRepositoryImpl(db)
	brokerRepo := broker.NewRabbitRepositoryImpl(cfg)
	userService := service.NewUserService(userRepository, brokerRepo)
	srv := httpServe.NewGRPCServer(userService, cfg)

	api.RegisterAuthServiceServer(grpcServer, srv)

	go monitoring.StartMetricsServer(cfg.MetricsPort)

	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	grpcErr := make(chan error, 1)
	go func() {
		log.Println("gRPC server starting on :8081")
		if err := grpcServer.Serve(listener); err != nil {
			grpcErr <- err
		}
	}()

	select {
	case err = <-grpcErr:
		log.Fatalf("gRPC server error: %v", err)
	case <-ctx.Done():
		log.Println("Received shutdown signal...")

		_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Println("Stopping gRPC server gracefully...")
		grpcServer.GracefulStop()

		log.Println("Closing database connections...")
		if err := db.Close(); err != nil {
			log.Printf("Database close error: %v", err)
		}

		log.Println("Server stopped gracefully")
	}
}
