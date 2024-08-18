package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	api_v1 "github.com/bryopsida/go-grpc-server-template/api/v1"
	"github.com/bryopsida/go-grpc-server-template/config"
	"github.com/bryopsida/go-grpc-server-template/datastore"
	"github.com/bryopsida/go-grpc-server-template/repositories/number"
	"github.com/bryopsida/go-grpc-server-template/services/increment"
	"google.golang.org/grpc"
)

func serveGrpc(server *grpc.Server, lis net.Listener) {
	slog.Info("Listening on :50051")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func runGrpc(ctx context.Context, server *grpc.Server, lis net.Listener) {
	go serveGrpc(server, lis)
	// wait for cancel signal
	<-ctx.Done()
	// stop the server
	slog.Info("Shutting down gRPC server...")
	server.GracefulStop()
}

func main() {
	slog.Info("Starting")
	config := config.NewViperConfig()
	slog.Info("Getting database")
	db, err := datastore.GetDatabase(config)
	if err != nil {
		log.Fatalf("failed to get database: %v", err)
	}
	defer db.Close()

	slog.Info("Getting number repository")
	repo := number.NewBadgerNumberRepository(db)

	slog.Info("Getting increment service")
	service := increment.NewIncrementService(repo, "counter")
	slog.Info("Creating gRPC server")
	// Create a new gRPC server
	server := grpc.NewServer()

	// Register the IncrementService
	api_v1.RegisterIncrementServiceServer(server, service)

	// Listen on a port
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	// ensure this is always called on func exit
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Run the server in a goroutine
	go runGrpc(ctx, server, lis)

	// Wait for a signal
	sig := <-sigChan
	slog.Info("Received signal", "signal", sig)
	// Cancel the context
	cancel()
	slog.Info("Server stopped")
}
