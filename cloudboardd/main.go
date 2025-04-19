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
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	listenAddr := flag.String("listen", ":50055", "cloudboard service listen address")
	attestAddr := flag.String("attest-addr", ":50051", "attestation service address")
	jobauthAddr := flag.String("jobauth-addr", ":50054", "job authorization service address")
	jobhelperBin := flag.String("jobhelper-bin", "./jobhelperd", "path to jobhelper binary")
	flag.Parse()
	fmt.Println("done config")
	// Connect to Attestation service
	attestConn, err := grpc.NewClient(*attestAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial attestation service: %v\n", err)
		os.Exit(1)
	}
	defer attestConn.Close()
	attestClient := attestation.NewAttestationClient(attestConn)

	// Connect to JobAuth service
	jobauthConn, err := grpc.NewClient(*jobauthAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to dial jobauth service: %v\n", err)
		os.Exit(1)
	}
	defer jobauthConn.Close()
	jobauthClient := jobauth.NewJobAuthClient(jobauthConn)

	// Note: JobHelper processes are spawned per-request using jobhelper binary

	// Start CloudBoard gRPC server
	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *listenAddr, err)
		os.Exit(1)
	}
	grpcServer := grpc.NewServer()
	srv := &server{
		attest:       attestClient,
		jobauth:      jobauthClient,
		jobhelperBin: *jobhelperBin,
		jobauthAddr:  *jobauthAddr,
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
type server struct {
	controller.UnimplementedCloudBoardServer
	attest       attestation.AttestationClient
	jobauth      jobauth.JobAuthClient
	jobhelperBin string
	jobauthAddr  string
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
// InvokeWorkload initiates a workload: first receives job metadata, obtains a token, then proxies streams.
func (s *server) InvokeWorkload(stream controller.CloudBoard_InvokeWorkloadServer) error {
	ctx := stream.Context()

	// 1. Receive initial job metadata from client
	initReq, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("failed to receive initial metadata: %w", err)
	}
	jobMetadata := initReq.Payload

	// 2. Generate a per-job token using the metadata
	tokenRes, err := s.jobauth.GenerateToken(ctx, &jobauth.GenerateTokenRequest{JobMetadata: jobMetadata})
	if err != nil {
		return fmt.Errorf("GenerateToken failed: %w", err)
	}

	// 3. Spawn JobHelper process, open stream and send the token
	helperStream, cleanup, err := spawnJobHelper(ctx, s.jobhelperBin, s.jobauthAddr)
	if err != nil {
		return fmt.Errorf("spawn jobhelper failed: %w", err)
	}
	defer cleanup()
	if err := helperStream.Send(&jobhelper.WorkloadRequest{Payload: tokenRes.Token}); err != nil {
		return fmt.Errorf("sending token to jobhelper failed: %w", err)
	}

	// 4. Proxy streams between client and JobHelper
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
