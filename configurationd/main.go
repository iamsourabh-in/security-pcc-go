package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/configuration"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("listen", ":50052", "configuration service listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *addr, err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	configuration.RegisterConfigurationServer(grpcServer, &server{})

	go func() {
		fmt.Printf("Configuration service listening on %s\n", *addr)
		if err := grpcServer.Serve(lis); err != nil {
			fmt.Fprintf(os.Stderr, "failed to serve: %v\n", err)
			os.Exit(1)
		}
	}()

	// Wait for termination signal.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	grpcServer.GracefulStop()
	fmt.Println("Configuration service shut down")
}

// server implements the Configuration gRPC service.
type server struct {
	configuration.UnimplementedConfigurationServer
}

// Register handles a client registration request.
// Register handles a client registration request by returning the current configuration.
func (s *server) Register(ctx context.Context, req *configuration.RegisterRequest) (*configuration.ConfigurationUpdate, error) {
   fmt.Printf("Configuration Register called with payload: %x\n", req.Registration)
   // Return a default configuration update (static payload for now)
   return &configuration.ConfigurationUpdate{UpdatePayload: []byte("default-config-v1")}, nil
}

// SuccessfullyAppliedConfiguration handles success acknowledgments.
// SuccessfullyAppliedConfiguration records a successful configuration application.
func (s *server) SuccessfullyAppliedConfiguration(ctx context.Context, req *configuration.ApplySuccessRequest) (*configuration.EmptyResponse, error) {
   fmt.Printf("Configuration applied successfully: %x\n", req.SuccessPayload)
   return &configuration.EmptyResponse{}, nil
}

// FailedToApplyConfiguration handles failure acknowledgments.
// FailedToApplyConfiguration records a failed configuration application.
func (s *server) FailedToApplyConfiguration(ctx context.Context, req *configuration.ApplyFailureRequest) (*configuration.EmptyResponse, error) {
   fmt.Printf("Configuration apply failed: %x\n", req.FailurePayload)
   return &configuration.EmptyResponse{}, nil
}

// CurrentConfigurationVersionInfo returns the version info.
// CurrentConfigurationVersionInfo returns info about the current configuration version.
func (s *server) CurrentConfigurationVersionInfo(ctx context.Context, req *configuration.EmptyRequest) (*configuration.VersionInfoResponse, error) {
   // Return static version info for now
   return &configuration.VersionInfoResponse{VersionInfo: []byte("version-1")}, nil
}
