package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/attestation"
	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/controller"
	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobauth"
	"github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobhelper"
	"google.golang.org/grpc"
)

func main() {
	listenAddr := flag.String("listen", ":50055", "cloudboard service listen address")
	attestAddr := flag.String("attest-addr", ":50051", "attestation service address")
	jobauthAddr := flag.String("jobauth-addr", ":50054", "job authorization service address")
	jobhelperAddr := flag.String("jobhelper-addr", ":50053", "job helper service address")
	flag.Parse()

	// Connect to Attestation service
	attestConn, err := grpc.Dial(*attestAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial attestation service: %v\n", err)
		os.Exit(1)
	}
	defer attestConn.Close()
	attestClient := attestation.NewAttestationClient(attestConn)

	// Connect to JobAuth service
	jobauthConn, err := grpc.Dial(*jobauthAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial jobauth service: %v\n", err)
		os.Exit(1)
	}
	defer jobauthConn.Close()
	jobauthClient := jobauth.NewJobAuthClient(jobauthConn)

	// Connect to JobHelper service
	jobhelperConn, err := grpc.Dial(*jobhelperAddr, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial jobhelper service: %v\n", err)
		os.Exit(1)
	}
	defer jobhelperConn.Close()
	jobhelperClient := jobhelper.NewJobHelperClient(jobhelperConn)

	// Start CloudBoard gRPC server
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *listenAddr, err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	srv := &server{
		attest:    attestClient,
		jobauth:   jobauthClient,
		jobhelper: jobhelperClient,
	}
	controller.RegisterCloudBoardServer(grpcServer, srv)

	go func() {
		fmt.Printf("CloudBoard service listening on %s\n", *listenAddr)
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
	fmt.Println("CloudBoard service shut down")
}

// server implements the CloudBoard gRPC service.
// server implements the CloudBoard gRPC service.
type server struct {
	controller.UnimplementedCloudBoardServer
	attest    attestation.AttestationClient
	jobauth   jobauth.JobAuthClient
	jobhelper jobhelper.JobHelperClient
}

// FetchAttestation retrieves a fresh attestation bundle.
func (s *server) FetchAttestation(ctx context.Context, req *controller.FetchAttestationRequest) (*controller.FetchAttestationResponse, error) {
	res, err := s.attest.Attest(ctx, &attestation.AttestRequest{})
	if err != nil {
		return nil, err
	}
	return &controller.FetchAttestationResponse{Bundle: res.Bundle}, nil
}

// InvokeWorkload initiates a workload by obtaining a token, forwarding the token to JobHelper,
// and proxying subsequent request/response streams.
func (s *server) InvokeWorkload(stream controller.CloudBoard_InvokeWorkloadServer) error {
	ctx := stream.Context()

	// 1. Generate a per-job token
	tokenRes, err := s.jobauth.GenerateToken(ctx, &jobauth.GenerateTokenRequest{JobMetadata: []byte("job-metadata")})
	if err != nil {
		return fmt.Errorf("GenerateToken failed: %w", err)
	}

	// 2. Open JobHelper stream
	helperStream, err := s.jobhelper.InvokeWorkload(ctx)
	if err != nil {
		return fmt.Errorf("InvokeWorkload (jobhelper) failed: %w", err)
	}
	// Send the token as the first message
	if err := helperStream.Send(&jobhelper.WorkloadRequest{Payload: tokenRes.Token}); err != nil {
		return fmt.Errorf("sending token to jobhelper failed: %w", err)
	}

	// 3. Proxy streams between client and JobHelper
	errCh := make(chan error, 2)

	// a) Client -> JobHelper
	go func() {
		for {
			req, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			if err := helperStream.Send(&jobhelper.WorkloadRequest{Payload: req.Payload}); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// b) JobHelper -> Client
	go func() {
		for {
			resp, err := helperStream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			if err := stream.Send(&controller.InvokeWorkloadResponse{Payload: resp.Payload}); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Wait for error or completion
	err = <-errCh
	helperStream.CloseSend()
	return err
}
