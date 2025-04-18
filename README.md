# Private Cloud Compute

This source code accompanies the Private Cloud Compute (PCC) [Security Guide](https://security.apple.com/documentation/private-cloud-compute/).

## About The Source Code

The source code in this repository includes components of PCC that implement security mechanisms and apply privacy policies. We provide this code to allow researchers and interested individuals to independently verify PCC's security and privacy characteristics and functionality.

## Getting Started

### Clone this repository

```bash
git clone https://github.com/apple/security-pcc.git
```

### Generate Go code for Go clients & servers

This repository includes Protocol Buffer definitions under `proto/`. To generate Go types and gRPC stubs:

- Install `protoc` (Protocol Buffer compiler)
- Install the Go plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```
- Ensure `$GOPATH/bin` is in your `PATH` so that `protoc` can find `protoc-gen-go` and `protoc-gen-go-grpc`.
- Run the generator script:
  ```bash
  scripts/gen-proto.sh
  ```

### Running Services

You can run each Go-based service in separate terminals. Ensure protoc-generated code is present.

```bash
# Start the JobAuth service (port 50054)
go run jobauthd/main.go
# Start the JobHelper service (port 50053)
go run jobhelper/main.go --jobauth-addr localhost:50054
# Start the Attestation service (port 50051)
go run attestationd/main.go
# Start the CloudBoard service (port 50055)
go run cloudboardd/main.go \
    --attest-addr localhost:50051 \
    --jobauth-addr localhost:50054 \
    --jobhelper-addr localhost:50053
```

### Contributions

The publication of this code is intended for security research and verification purposes only. Please see [CONTRIBUTING.md](CONTRIBUTING.md) for more information.

## Contact

To report security issues with this code, please use the instructions [at this page](https://support.apple.com/en-us/102549).
