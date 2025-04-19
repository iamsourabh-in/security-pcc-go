# Private Cloud Compute ☁️

This source code accompanies the Private Cloud Compute (PCC) [Security Guide](https://security.apple.com/documentation/private-cloud-compute/).

---


## About The Source Code 🛡️

The source code in this repository includes components of PCC that implement security mechanisms and apply privacy policies. We provide this code to allow researchers and interested individuals to independently verify PCC's security and privacy characteristics and functionality.

This repository represents Apple’s “Private Cloud Compute” monorepo. It comprises several related but independently buildable Swift/C/C++ components implementing PCC’s security, attestation, orchestration, metrics, diagnostics, OS‑level plumbing, developer tools, and a “Thimble” client/daemon.

---

## Repository Overview 🗺️

Here’s a high‑level look at the main components:

*   **AppleComputeEnsembler**
    *   A Swift‐based “ensemble” service including daemons (`ensembled`), a CLI (`ensemblectl`), and a framework/DSL for orchestrating multi‑node compute tasks.
*   **CloudAttestation** 🔒
    *   A Swift library handling attestation assets, validators, policy engine, provisioning certificates, transparency logs, HPKE channels, etc.
*   **CloudBoard** ⚙️
    *   A comprehensive Swift/XPC/gRPC‑based orchestration toolkit. Includes:
        *   Daemons (`cb_*`), XPC & gRPC APIs.
        *   Controllers & frameworks for attestation, configuration, job scheduling, identity, health, preferences, logging, metrics, etc.
*   **CloudMetrics** 📊
    *   Swift frameworks & daemons for collecting, filtering, and exporting metrics (OpenTelemetry, Splunk, CoreAnalytics), plus test applications.
*   **CloudRemoteDiagnostics** 🩺
    *   CLI, core libraries (C++ ping/tcpdump/stats), daemon, and framework for running remote diagnostics commands over XPC.
*   **darwinOSBits** 🍎
    *   Low‑level system services and tools:
        *   `darwin‑init`: System initialization.
        *   `secureconfig`/`secureconfigd`/`ctl`: Secure configuration management.
        *   `SecureConfigDB`: Database for secure configuration.
        *   `splunkloggingd`: System log offload daemon.
*   **SecurityMonitorLite** 👀
    *   A lightweight daemon monitoring `exec`, network, SSH login/logout events and exporting them (e.g., to Splunk).
*   **srd_tools** 🛠️
    *   Security‑research developer tools (Swift) for Mach‑O introspection, “VRE” (virtualized runtimes) CLI, VM helpers, etc., with unit tests.
*   **Thimble**  CLIENT/DAEMON
    *   A Swift “trusted cloud compute” client & daemon implementing authenticated/encrypted request framing over XPC, token management, rate‑limiting, server‑driven config, telemetry, plus a supporting framework and XPC protocol.

---

## CloudBoard Deep Dive 🔬

CloudBoard is Apple’s Swift‑based orchestration layer within PCC. It's an ecosystem of daemons, APIs, libraries, and frameworks responsible for node attestation, configuration management, job scheduling, identity, health, logging/metrics, and inter‑node coordination.

### 1. Core Daemons & Responsibilities

*   `cloudboardd`: The main controller daemon on each node. Manages configuration, interacts with helper daemons, runs XPC/gRPC servers, enforces policies, schedules workloads, and reports health/metrics.
*   `cb_attestationd`: A Swift/XPC front-end for the `CloudAttestation` library. Used by `cloudboardd` for node attestation, fetching proofs, and validating certificates/policies.
*   `cb_configurationd`: Manages node configuration ("what should I be doing?"). Implements a registry state machine, fetches manifests/configs, applies changes, and writes local state.
*   `cb_jobhelper`: A less-privileged helper daemon. `cloudboardd` delegates workload launch and I/O management to it.
*   `cloudboardd_follower`: A standby/cross‑site follower mirroring registry state from a leader `cloudboardd` for high availability.
*   `cb_jobauthd`: A small daemon responsible for signing or validating per‑job tokens/credentials to authorize workload execution.

### 2. API & IPC Layers 🔌

*   **XPC**: Each daemon exposes an XPC interface (defined under `CloudBoard*API/XPC`) for local inter-process communication.
*   **gRPC/Protobuf**: `cloudboardd` hosts a gRPC server for remote management (e.g., via a control plane). Protobuf definitions are under `CloudBoardDCore/Models` and `Resources`.
*   **CloudBoardAsyncXPC**: A shared Swift helper library for robust XPC communication (handles connection back-off, serialization, error mapping).

### 3. Shared Frameworks & Libraries 📚

*   **CloudBoardCommon**: Provides utilities like Futures/Promises, retry/backoff logic, `CFPreferences` encoding, `launchd` helpers, secure-config plumbing, node info access, etc.
*   **CloudBoardPlatformUtilities**: Includes service discovery (mDNS), reliable connections, watchdog services, metrics exporters, job quiescence monitoring, and lifecycle management.
*   **CloudBoardLogging & CloudBoardMetrics**: Facades for structured logging and metrics collection, integrated with OSLog and OpenTelemetry.
*   **CloudBoardController Framework**: A high-level API for clients to request task execution and track progress.

### 4. Configuration Subsystem (`CloudBoardConfigurationDCore`) 🚦

*   Implements a state machine (see `Docs/`) that reads local registry, polls/subscribes to remote configuration services, computes deltas, updates the registry, and notifies `cloudboardd` of changes.

### 5. Attestation Subsystem (`CloudBoardAttestationDCore`) ✅

*   Utilizes the `CloudAttestation` library to gather hardware/cryptographic identity, build/validate attestation bundles, enforce local policies, and produce an attested node identity.

### 6. Job Scheduling & Execution ▶️

*   `cloudboardd` decides when to run workloads.
*   It generates job tokens via `cb_jobauthd`.
*   It invokes `cb_jobhelper` via XPC to launch processes with appropriate sandboxing/entitlements.
*   It uses `CloudBoardPlatformUtilities` for health monitoring.
*   It streams logs & metrics via its API or XPC.

### 7. Clients & Test Harnesses 🧪

*   **LocalCloudBoardClient**: A lightweight Swift client for interacting with the local daemon.
*   **null‑cloud‑controller-cli, NullCloudApp, NullCloudController**: Mock controllers and applications for end-to-end testing without a real control plane.

### 8. Observability & Health ❤️‍🩹

*   Components emit structured logs (`CloudBoardLogging`) and metrics (`CloudBoardMetrics`).
*   Health check endpoints (XPC/gRPC) allow external monitoring.

### 9. How It All Fits Together ✨

1.  `launchd` starts the various CloudBoard daemons.
2.  On boot, `cb_configurationd` fetches the node's registry.
3.  `cloudboardd` reads the registry, uses `cb_attestationd` to prove identity, publishes its APIs, and awaits workloads.
4.  A Control Plane (via gRPC/XPC) can issue commands like "run workload X" or "update configuration".
5.  `cloudboardd` enforces policy, uses `cb_jobauthd` for tokens, and `cb_jobhelper` to launch the workload.
6.  In Follower Mode, `cloudboardd_follower` mirrors state from a leader for redundancy or load balancing.

**In short:** CloudBoard acts as the comprehensive "private-cloud agent" on each device, managing secure identity, configuration, workload execution, and observability through modular Swift components and standard communication protocols.

---

## Technical Details 💻

*   **Languages**: Primarily Swift (with some C/Objective-C bridging), with a few C/C++ modules (`secureconfig`, remote diagnostics core).
*   **Build System**: Components are typically structured as Swift Packages or Xcode projects. CLI tools often use `Swift ArgumentParser`. IPC relies heavily on XPC, while RPC uses gRPC/Protobuf.
*   **Build Process**: There is no single "root" build command. Open the Xcode workspace (`.xcworkspace`) or build individual Swift packages/targets as needed.
*   **Purpose**: This repository accompanies Apple’s Private Cloud Compute Security Guide, offering transparency into the security mechanisms of the PCC stack.
*   **Contribution Model**: This is a read‑only repository for external contributors, focused on enabling security research and verification.

---

## Exploring the Codebase 🧭

*   Start with this top‑level `README.md`.
*   Each subfolder typically contains its own `Package.swift` or `.xcodeproj` file, often accompanied by a specific README or documentation comments within the code.
*   The `CloudAttestation` and `AppleComputeEnsembler` modules are central to PCC’s cryptographic protocols.
*   `CloudBoard` orchestrates these components.
*   `darwinOSBits` and `SecurityMonitorLite` operate at the OS interface level.

Feel free to delve into specific components!

---

## Getting Started (Go Clients/Servers) 🚀

This repository includes Protocol Buffer definitions under `proto/` which can be used to generate Go client and server code. *Note: This is for interacting with potential Go-based services using the defined gRPC interfaces, not for building the core Swift/C++ PCC components.*

### Prerequisites

1.  Install `protoc` (Protocol Buffer compiler).
2.  Install the Go plugins for `protoc`:
    ```bash
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    ```
3.  Ensure your `$GOPATH/bin` directory is included in your system's `PATH` environment variable so `protoc` can locate the installed Go plugins (`protoc-gen-go` and `protoc-gen-go-grpc`).

### Generate Go Code

Run the provided generator script:

```bash
scripts/gen-proto.sh

go list -mod=readonly -m -u -e -json all

# Start the example JobAuth service (listens on port 50054)
go run jobauthd/main.go

# Start the example JobHelper service (listens on port 50053)
# Connects to JobAuth service at localhost:50054
go run jobhelper/main.go --jobauth-addr localhost:50054

# Start the example Attestation service (listens on port 50051)
go run attestationd/main.go

# Start the example CloudBoard service (listens on port 50055)
# Connects to Attestation, JobAuth, and JobHelper services
go run cloudboardd/main.go \
    --attest-addr localhost:50051 \
    --jobauth-addr localhost:50054 \
    --jobhelper-addr localhost:50053

```
