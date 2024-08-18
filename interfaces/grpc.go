package interfaces

import "net"

// GrpcServer is an interface for gRPC server
type GrpcServer interface {
	// Serve starts the gRPC server
	// - net.Listener: the listener to bind the server to
	// returns an error if the server fails to start
	Serve(net.Listener) error
	// GracefulStop stops the gRPC server gracefully
	GracefulStop()
}
