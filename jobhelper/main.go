package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobauth"
	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobhelper"
	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("listen", ":50053", "job helper service listen address")
	jobauthAddr := flag.String("jobauth-addr", ":50054", "job authorization service address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *addr, err)
		os.Exit(1)
	}

	// Connect to JobAuth service for token validation
	jobauthConn, err := grpc.Dial(*jobauthAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial jobauth service: %v\n", err)
		os.Exit(1)
	}
	defer jobauthConn.Close()
	jobauthClient := jobauth.NewJobAuthClient(jobauthConn)

	grpcServer := grpc.NewServer()
	// Register server with JobAuth client
	jobhelper.RegisterJobHelperServer(grpcServer, &server{jobauth: jobauthClient})

	go func() {
		fmt.Printf("JobHelper service listening on %s\n", *addr)
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
	fmt.Println("JobHelper service shut down")
}

// server implements the JobHelper gRPC service with token validation.
type server struct {
	jobhelper.UnimplementedJobHelperServer
	jobauth jobauth.JobAuthClient
}

// InvokeWorkload handles streaming workload requests and responses.
// The first message must contain a valid token.
func (s *server) InvokeWorkload(stream jobhelper.JobHelper_InvokeWorkloadServer) error {
	ctx := stream.Context()
	// 1) Receive token message
	tokenReq, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive token: %w", err)
	}
	// Validate token
	valRes, err := s.jobauth.ValidateToken(ctx, &jobauth.ValidateTokenRequest{Token: tokenReq.Payload})
	if err != nil {
		return fmt.Errorf("ValidateToken error: %w", err)
	}
	if !valRes.Valid {
		return fmt.Errorf("invalid token provided")
	}
	// 2) Process subsequent workload requests
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Println("Received workload request:", string(req.Payload))
		// For demonstration, simply echo the payload
		resp := &jobhelper.WorkloadResponse{Payload: []byte(base64.StdEncoding.EncodeToString(req.Payload) + " !!")}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}

// Teardown handles teardown of workloads.
func (s *server) Teardown(ctx context.Context, req *jobhelper.TeardownRequest) (*jobhelper.EmptyResponse, error) {
	// Log teardown request payload
	fmt.Printf("JobHelper Teardown called with payload: %x\n", req.Payload)
	// Teardown logic can be added here (cleanup resources)
	return &jobhelper.EmptyResponse{}, nil
}
