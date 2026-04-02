# CLAUDE.md

This file provides guidance for the protobuf code generator **genpb**.

## Project Overview

Generates **Go** (server) and **Rust / Godot** (client) code from `.proto` files.  
**C# generation was removed** (`gen_cs.go` deleted); do not reference `--lang Pb` or `--cs_out`.

Rust output is driven by a **`ProtocolManifest`** built from `FileDescriptorSet` (`protocol.desc`), not regex on `.proto` sources.

## Commands

### Build the tool

```bash
go build -buildvcs=false -o genpb      # Unix
go build -buildvcs=false -o genpb.exe  # Windows
```

If the repo has mixed VCS roots, always pass `-buildvcs=false` to `go build` / `go run`.

### Scripts (`gen.bat` / `gen.sh`)

Wrapper scripts match `main.go`. Args: `flag proto_in go_out lang tools_dir [rust_out] [godot_out]`.  
**No** `--cs_out`, **no** `Pb` language.

Typical Rust + Godot client generation:

```bash
go run -buildvcs=false . --lang rust --flag client \
  --rust_out  ../../gclient/rust/lib/gnet/src/gen \
  --godot_out ../../gclient/rust/gdbridge/src/gen
```

Go server only:

```bash
go run -buildvcs=false . --lang go --flag server --go_out ./pb
```

### Command-line options

| Flag | Default | Description |
|------|---------|-------------|
| `--lang` | `all` | `go`, `rust`, or `all` |
| `--proto_in` | `./proto` | Proto input directory |
| `--go_out` | `./pb` | Go output directory |
| `--rust_out` | _(empty)_ | Rust outputs: `pb.rs`, `typed_protocol.rs`, `protocol_manifest.json`, `protocol.desc` |
| `--godot_out` | _(empty)_ | Writes `godot_bridge_gen.rs` (optional; use with `--rust_out`) |
| `--flag` | `server` | `server` (all protos) or `client` (excludes `data_srv.proto`, `data_fwd.proto`) |
| `--tools_dir` | `../proto` | Directory with `protoc`, `protoc-gen-go`, `protoc-gen-prost` |

`--lang rust` requires `--rust_out`.

## File structure

```
genpb/
├── main.go       # Entry + config
├── manifest.go   # ProtocolManifest + descriptor parsing + recursive fingerprints
├── gen_go.go     # Go generation (unchanged contract for server)
├── gen_rust.go   # prost + manifest + typed_protocol.rs + optional godot_bridge_gen.rs
├── gen.bat
├── gen.sh
├── CLAUDE.md
├── proto/
└── pb/           # Default Go output (local)
```

## Generated files

### Go (server)

| Output | Purpose |
|--------|---------|
| `*.pb.go` | Messages / enums from protoc |
| `cmd.ext.go` | `Package`, `Unmarshal`, per-message `Key()` / `Marshal()` |

### Rust (`--rust_out`)

| File | Purpose |
|------|---------|
| `pb.rs` | Prost types from `protoc-gen-prost` |
| `typed_protocol.rs` | `EKey`, `ClientMessage`, `ServerMessage`, encode/decode, `COMPILED_FINGERPRINTS` |
| `protocol_manifest.json` | Full manifest (replaces legacy `protocol_meta.json`) |
| `protocol.desc` | Binary `FileDescriptorSet` for `prost-reflect` / hotfix |

### Godot (`--godot_out`)

| File | Purpose |
|------|---------|
| `godot_bridge_gen.rs` | `*Gd` classes, `NetEventGd`, `server_message_to_event`, `hotfix_to_event` |

## Go `cmd.ext.go` pattern (reference)

```go
var _parser = NewParser()
func init() { _parser.Load() }

type Package struct { ... }
func NewPackage(msg proto.Message, errs ...EErrorCode_T) *Package
func (p *Package) Key() EKey_T
func (p *Package) Marshal() ([]byte, error)
func Unmarshal(key EKey_T, data []byte) proto.Message
```

## EKey / proto layout

Defined in `proto/cmd.proto` (`message EKey { enum T { ... } }`).  
`manifest.go` reads enum values from the descriptor for `cmd.proto`.  
Request/response/dispatch messages live in `cmd_req.proto`, `cmd_rsp.proto`, `cmd_dsp.proto`; direction is inferred from file + naming (`Req*` / `Rsp*` / `Dsp*`).

## Related paths

- `gclient/rust/lib/gnet` – consumes `gen/pb.rs` + `gen/typed_protocol.rs`
- `gclient/rust/gdbridge` – `include!("gen/godot_bridge_gen.rs")`
- Server Go code – generated under configured `--go_out`
