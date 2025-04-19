package main

import (
   "context"
   "crypto/rand"
   "encoding/hex"
   "encoding/json"
   "flag"
   "fmt"
   "net"
   "os"
   "os/signal"
   "syscall"

   "github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/attestation"
   "google.golang.org/grpc"
)

func main() {
	addr := flag.String("listen", ":50051", "attestation service listen address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to listen on %s: %v\n", *addr, err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	attestation.RegisterAttestationServer(grpcServer, &server{})

	go func() {
		fmt.Printf("Attestation service listening on %s\n", *addr)
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
	fmt.Println("Attestation service shut down")
}

// server implements the Attestation gRPC service.
type server struct {
	attestation.UnimplementedAttestationServer
}

// Attest handles an attestation request by returning a simple attestation bundle.
func (s *server) Attest(ctx context.Context, req *attestation.AttestRequest) (*attestation.AttestResponse, error) {
   // Gather basic host information
   host, err := os.Hostname()
   if err != nil {
       return nil, fmt.Errorf("failed to get hostname: %w", err)
   }
   // Generate a random nonce
   nonce := make([]byte, 16)
   if _, err := rand.Read(nonce); err != nil {
       return nil, fmt.Errorf("failed to generate nonce: %w", err)
   }
   // Prepare attestation bundle
   bundle := map[string]string{
       "hostname": host,
       "nonce":    hex.EncodeToString(nonce),
   }
   data, err := json.Marshal(bundle)
   if err != nil {
       return nil, fmt.Errorf("failed to marshal bundle: %w", err)
   }
   return &attestation.AttestResponse{Bundle: data}, nil
}
