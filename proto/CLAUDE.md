# CLAUDE.md

This file provides guidance for the protobuf source definitions in this directory.

## Project Overview

Protocol buffer (protobuf) definitions for client–server communication. **Do not run raw `protoc` for game pipelines** unless you know what you are doing: the supported entry point is **`../` (genpb)**, which generates Go (server) and Rust/Godot (client) and builds **`ProtocolManifest`** from descriptors.

**genpb no longer generates C#**; legacy `protoc --csharp_out` snippets below are not part of the maintained workflow.

## Tool Versions

- protoc: 3.20.x+ (project-tested)
- protoc-gen-go: v1.36+
- protoc-gen-prost: install via `cargo install protoc-gen-prost` (Rust client)

## Proto Files

| File | Purpose | Key Definitions |
|------|---------|-----------------|
| `enum.proto` | Enumerations | EErrorCode, EKickType, … |
| `cmd.proto` | Command keys | `message EKey { enum T { … } }` |
| `cmd_req.proto` | Client → server | `Req*` messages |
| `cmd_rsp.proto` | Server → client (RPC 响应) | `Rsp*` messages |
| `cmd_dsp.proto` | Server push | `Dsp*` messages |
| `data.proto` | Shared structs | RoleSummaryData, LoginData, Vector, … |
| `data_fwd.proto` | Forward messages | (server export) |
| `data_srv.proto` | Server-only data | (server export) |

## Regenerating Code

From **`comm/tools/genpb`** (parent of `proto/`):

```bash
# Server Go
go run -buildvcs=false . --lang go --flag server --go_out <path/to/server/pb>

# Client Rust + optional Godot bridge
go run -buildvcs=false . --lang rust --flag client \
  --rust_out  ../../gclient/rust/lib/gnet/src/gen \
  --godot_out ../../gclient/rust/gdbridge/src/gen
```

See **`../README.md`** for full flags and outputs (`pb.rs`, `typed_protocol.rs`, `protocol_manifest.json`, `protocol.desc`, `godot_bridge_gen.rs`).

## Related Paths

- `../` — genpb tool (`main.go`, `manifest.go`, `gen_rust.go`, …)
- `../../gclient/rust/lib/gnet` — includes generated `gen/pb.rs` and `gen/typed_protocol.rs`
- `../../gclient/rust/gdbridge` — includes generated `gen/godot_bridge_gen.rs`
