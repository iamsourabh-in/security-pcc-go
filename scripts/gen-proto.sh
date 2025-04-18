#!/usr/bin/env bash
set -e

# Ensure Go bin directory is in PATH for protoc plugins
GO_BIN=$(go env GOBIN)
if [ -z "$GO_BIN" ]; then
  GO_BIN=$(go env GOPATH)/bin
fi
export PATH="$GO_BIN:$PATH"
# Check for protoc
# Check for protoc installation
if ! command -v protoc >/dev/null 2>&1; then
  echo "Error: protoc not found. Install protoc to continue." >&2
  exit 1
fi
# Ensure protoc binary is executable
PROTOC_BIN=$(command -v protoc)
if [ ! -x "$PROTOC_BIN" ]; then
  echo "Error: protoc binary at '$PROTOC_BIN' is not executable. Permission denied." >&2
  echo "Please run: chmod +x '$PROTOC_BIN' or reinstall protoc (e.g., via Homebrew or your package manager)." >&2
  exit 1
fi

# Check for Go protobuf plugin
if ! command -v protoc-gen-go >/dev/null 2>&1; then
  echo "Error: protoc-gen-go not found. Install it with:" >&2
  echo "  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest" >&2
  exit 1
fi

# Check for Go gRPC plugin
if ! command -v protoc-gen-go-grpc >/dev/null 2>&1; then
  echo "Error: protoc-gen-go-grpc not found. Install it with:" >&2
  echo "  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest" >&2
  exit 1
fi

echo "Generating Go code from .proto files..."
for f in $(find proto -name '*.proto'); do
  protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative "$f"
done
echo "Done."