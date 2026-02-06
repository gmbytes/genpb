# CLAUDE.md

This file provides guidance for the protobuf code generator.

## Project Overview

Generates Go and C# code from `.proto` files for game server client-server communication.

## Commands

```bash
# Generate all (Go + C#)
go run main.go

# Generate Go only
go run main.go --lang go

# Generate C# only
go run main.go --lang cs
```

## File Structure

```
genpb/
├── main.go      # Entry point + config
├── gen_go.go    # Go code generation
├── gen_cs.go    # C# code generation
├── CLAUDE.md    # This file
└── pb/          # Generated files
```

## Generated Files

| Language | File | Purpose |
|----------|------|---------|
| Go | `pb/cmd.pb.go` | Command/Message ID definitions |
| Go | `pb/cmd.ext.go` | Parser + message helpers |
| Go | `pb/cmd_req.pb.go` | Request messages |
| Go | `pb/cmd_rsp.pb.go` | Response messages |
| Go | `pb/cmd_dsp.pb.go` | Dispatch messages |
| Go | `pb/enum.pb.go` | Enumeration types |
| C# | `pb/Cmd.cs` | Command keys + helpers |

## Generated Code Features

### Go (`cmd.ext.go`)
```go
// Global parser with auto-init
var _parser = NewParser()

func init() {
    _parser.Load()
}

func Unmarshal(key EKey_T, data []byte) proto.Message

// Message.Key() method
func (msg *ReqLogin) Key() pb.EKey
```

### C# (`Cmd.cs`)
```csharp
// Command keys enum
public enum EKey
{
    Login = 1,
    CreateRole = 2,
    // ...
}

// Message helper extension
public static class CmdExtensions
{
    public static Cmd.EKey GetKey(this ReqLogin msg);
}
```

## Message ID Ranges (EKey)

| Range | Category |
|-------|----------|
| 1-9 | Login flow |
| 10-19 | Heartbeat |
| 20-29 | Scene management |
| 40000+ | Server sync/push |

## Proto Naming Conventions

| Proto File | Message Prefix | Enum Mapping |
|------------|---------------|--------------|
| `cmd_req.proto` | `Req*` | `ReqLogin` → `EKey_Login` |
| `cmd_rsp.proto` | `Rsp*` | `RspLogin` → `EKey_Login` |
| `cmd_dsp.proto` | `RspSync*` | `RspSyncLoginFast` → `EKey_SyncLoginFast` |

## Related Projects

- `E:/tools/proto` - Protobuf source definitions
- `E:/tools/server` - Game server using generated types
- `E:/tools/client/sgame` - C# client
