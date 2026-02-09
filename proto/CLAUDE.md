# CLAUDE.md

This file provides guidance for the protobuf source definitions.

## Project Overview

Protocol buffer (protobuf) definitions for a multiplayer game's client-server communication. Proto3 syntax generates both Go and C# code.

## Tool Versions

- protoc: 3.20.3
- protoc-gen-go: v1.36.11+
- protoc-gen-go-grpc: v1.3+

## Proto Files

| File | Purpose | Key Definitions |
|------|---------|-----------------|
| `enum.proto` | Enumerations | EErrorCode, EKickType, ERoleType, EEntityType, ESceneType |
| `cmd.proto` | Command keys | EKey.T enum (message IDs) |
| `cmd_req.proto` | Client requests | ReqLogin, ReqCreateRole, ReqPing, ReqEnterScene |
| `cmd_rsp.proto` | Server responses | RspLogin, RspCreateRole, RspPing, RspEnterScene |
| `cmd_dsp.proto` | Server dispatches | RspSyncLoginFast, RspSyncLoginData, RspSyncKickRole |
| `data.proto` | Data structures | SRoleSummaryData, SLoginData, SObject, SItem, SAttr |
| `data_fwd.proto` | Forward messages | SFwdCheckDistance, SFwdKick |
| `data_srv.proto` | Server data | SOrderInfo (payment) |

## Message ID Ranges

| Range | Category |
|-------|----------|
| 1-9 | Login flow |
| 10-19 | Heartbeat |
| 20-29 | Scene management |
| 40000+ | Server sync/push |

## Regenerating Code

```bash
# Generate Go code
./gen_proto.sh  # Linux/Mac
gen_proto.bat   # Windows

# Or manually:
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       --csharp_out=gen/Pb/ \
       *.proto
```

## Output Directories

- `gen/pb/` - Generated Go code (imported by `E:/tools/genpb/`)
- `gen/cs/` - Generated C# code (imported by `E:/tools/client/sgame/`)

## Related Projects

- `E:/tools/genpb` - Generated Go protobuf library
- `E:/tools/client/sgame` - C# client using generated code
- `E:/tools/server` - Go server using generated code
