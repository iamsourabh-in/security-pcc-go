package main

import (
  "context"
  "flag"
  "fmt"
  "net"
  "os"
  "os/signal"
  "syscall"

  "github.com/apple/security-pcc-go/proto/cloudboard/configuration"
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
func (s *server) Register(ctx context.Context, req *configuration.RegisterRequest) (*configuration.ConfigurationUpdate, error) {
  // TODO: integrate registry state machine logic.
  return &configuration.ConfigurationUpdate{}, nil
}

// SuccessfullyAppliedConfiguration handles success acknowledgments.
func (s *server) SuccessfullyAppliedConfiguration(ctx context.Context, req *configuration.ApplySuccessRequest) (*configuration.EmptyResponse, error) {
  // TODO: record successful apply.
  return &configuration.EmptyResponse{}, nil
}

// FailedToApplyConfiguration handles failure acknowledgments.
func (s *server) FailedToApplyConfiguration(ctx context.Context, req *configuration.ApplyFailureRequest) (*configuration.EmptyResponse, error) {
  // TODO: record failed apply.
  return &configuration.EmptyResponse{}, nil
}

// CurrentConfigurationVersionInfo returns the version info.
func (s *server) CurrentConfigurationVersionInfo(ctx context.Context, req *configuration.EmptyRequest) (*configuration.VersionInfoResponse, error) {
  // TODO: return current version info.
  return &configuration.VersionInfoResponse{}, nil
}