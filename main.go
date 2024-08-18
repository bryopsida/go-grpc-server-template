package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	api_v1 "github.com/bryopsida/go-grpc-server-template/api/v1"
	"github.com/bryopsida/go-grpc-server-template/config"
	"github.com/bryopsida/go-grpc-server-template/datastore"
	"github.com/bryopsida/go-grpc-server-template/interfaces"
	"github.com/bryopsida/go-grpc-server-template/repositories/number"
	"github.com/bryopsida/go-grpc-server-template/services/increment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func serveGrpc(server *grpc.Server, lis net.Listener) {
	address := lis.Addr().String()
	slog.Info("Listening on ", "address", address)
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

func buildTlsCert(cert, key string) (*tls.Certificate, error) {
	// Convert PEM strings to byte slices
	certPEM := []byte(cert)
	keyPEM := []byte(key)

	// Parse the certificate and key
	tlsCertificate, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate and key: %v", err)
	}

	return &tlsCertificate, nil
}

func buildServerCredentials(config interfaces.IConfig) (credentials.TransportCredentials, error) {
	cert := config.GetServerCert()
	key := config.GetServerKey()
	ca := config.GetServerCA()
	if cert == "" || key == "" || ca == "" {
		return nil, fmt.Errorf("missing required TLS configuration")
	}
	tlsCertificate, err := buildTlsCert(cert, key)
	if err != nil {
		return nil, fmt.Errorf("failed to build TLS certificate: %v", err)
	}
	creds := credentials.NewServerTLSFromCert(tlsCertificate)
	return creds, nil
}

func buildGrpcOptions(config interfaces.IConfig) []grpc.ServerOption {
	options := []grpc.ServerOption{}
	if config.IsTlsEnabled() {
		creds, err := buildServerCredentials(config)
		if err != nil {
			log.Fatalf("failed to build server credentials: %v", err)
		}
		options = append(options, grpc.Creds(creds))
	}
	return options
}

func buildGrpcServer(options []grpc.ServerOption) *grpc.Server {
	return grpc.NewServer(options...)
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
	options := buildGrpcOptions(config)
	server := buildGrpcServer(options)

	// Register the IncrementService
	api_v1.RegisterIncrementServiceServer(server, service)

	// Listen on a port
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", config.GetServerAddress(), config.GetServerPort()))
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
