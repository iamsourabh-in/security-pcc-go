
### Clone this repository

```bash
git clone https://github.com/iamsourabh-in/security-pcc-go.git
```


### 1. Core Daemons & Responsibilities

*   `cloudboardd`: The main controller daemon on each node. Manages configuration, interacts with helper daemons, runs XPC/gRPC servers, enforces policies, schedules workloads, and reports health/metrics.
*   `cb_attestationd`: A Swift/XPC front-end for the `CloudAttestation` library. Used by `cloudboardd` for node attestation, fetching proofs, and validating certificates/policies.
*   `cb_configurationd`: Manages node configuration ("what should I be doing?"). Implements a registry state machine, fetches manifests/configs, applies changes, and writes local state.
*   `cb_jobhelper`: A less-privileged helper daemon. `cloudboardd` delegates workload launch and I/O management to it.
*   `cloudboardd_follower`: A standby/cross‑site follower mirroring registry state from a leader `cloudboardd` for high availability.
*   `cb_jobauthd`: A small daemon responsible for signing or validating per‑job tokens/credentials to authorize workload execution.

### IPC Layers 🔌

*   **UDS**: Each daemon exposes an UDS  interface (Unix Domain Socket) for local inter-process communication.
*   **gRPC/Protobuf**: `cloudboardd` hosts a gRPC server for remote management (e.g., via a control plane). Protobuf definitions are under `CloudBoardDCore/Models` and `Resources`.


### 6. Job Scheduling & Execution ▶️

*   `cloudboardd` decides when to run workloads.
*   It generates job tokens via `cb_jobauthd`.
*   It creates minimum set of workers `cb_jobhelper` and send payload via UDS to communicate processes with appropriate sandboxing/entitlements.
*   It uses `CloudBoardPlatformUtilities` for health monitoring.
*   It streams logs & metrics via its API or UDS