package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	helperpb "github.com/iamsourabh-in/security-pcc-go/proto/cloudboard/jobhelper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// spawnJobHelper starts a new JobHelper process listening on an ephemeral port,
// waits for it to be ready, dials it, and opens a workload stream.
// Returns the stream, a cleanup function to close and kill the process, and any error.
func spawnJobHelper(ctx context.Context, binPath, jobauthAddr string) (helperpb.JobHelper_InvokeWorkloadClient, func() error, error) {
	// Prepare command: listen on random port
	cmd := exec.CommandContext(ctx, binPath,
		"--listen", ":50053",
		"--jobauth-addr", jobauthAddr,
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("jobhelper stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, fmt.Errorf("jobhelper stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, fmt.Errorf("starting jobhelper: %w", err)
	}
	// Stream stderr for logging
	go io.Copy(os.Stderr, stderr)

	// Read stdout until listening line appears
	reader := bufio.NewReader(stdout)
	const prefix = "JobHelper service listening on "
	var addr string
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			cmd.Process.Kill()
			return nil, nil, fmt.Errorf("reading jobhelper stdout: %w", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			addr = strings.TrimPrefix(line, prefix)
			break
		}
	}
	if addr == "" {
		cmd.Process.Kill()
		return nil, nil, fmt.Errorf("timeout waiting for jobhelper listening message")
	}

	// Connect to the JobHelper gRPC server
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		cmd.Process.Kill()
		return nil, nil, fmt.Errorf("dial jobhelper at %s: %w", addr, err)
	}
	client := helperpb.NewJobHelperClient(conn)
	stream, err := client.InvokeWorkload(ctx)
	if err != nil {
		conn.Close()
		cmd.Process.Kill()
		return nil, nil, fmt.Errorf("jobhelper InvokeWorkload: %w", err)
	}

	// cleanup closes stream, connection, and kills the process
	cleanup := func() error {
		stream.CloseSend()
		conn.Close()
		if cmd.Process != nil {
			cmd.Process.Kill()
			_, _ = cmd.Process.Wait()
		}
		return nil
	}
	return stream, cleanup, nil
}
