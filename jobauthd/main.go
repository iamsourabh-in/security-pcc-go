package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobauth"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("listen", ":50054", "job auth service listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *addr, err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	// Initialize token store
	authSrv := &server{tokens: make(map[string]struct{})}
	jobauth.RegisterJobAuthServer(grpcServer, authSrv)

	go func() {
		fmt.Printf("JobAuth service listening on %s\n", *addr)
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
	fmt.Println("JobAuth service shut down")
}

// server implements the JobAuth gRPC service with an in-memory token store.
type server struct {
	jobauth.UnimplementedJobAuthServer
	mu     sync.Mutex
	tokens map[string]struct{}
}

// GenerateToken issues a new token for a job.
func (s *server) GenerateToken(ctx context.Context, req *jobauth.GenerateTokenRequest) (*jobauth.GenerateTokenResponse, error) {
	// Generate a random 32-byte token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}
	token := hex.EncodeToString(b)
	// Store the token
	s.mu.Lock()
	s.tokens[token] = struct{}{}
	s.mu.Unlock()
	return &jobauth.GenerateTokenResponse{Token: []byte(token)}, nil
}

// ValidateToken checks if a token is valid.
func (s *server) ValidateToken(ctx context.Context, req *jobauth.ValidateTokenRequest) (*jobauth.ValidateTokenResponse, error) {
	token := string(req.Token)
	s.mu.Lock()
	_, ok := s.tokens[token]
	s.mu.Unlock()
	return &jobauth.ValidateTokenResponse{Valid: ok}, nil
}
