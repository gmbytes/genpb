# CLAUDE.md

This file provides guidance for the protobuf code generator.

## Project Overview

Generates Go and C# code from `.proto` files for game server client-server communication.

## Commands

### Scripts (recommended)

脚本会调用**已编译的** `genpb.exe`（Windows）或 `genpb`（Linux/macOS），优先使用脚本同目录下的可执行文件。需先编译：`go build -buildvcs=false -o genpb`（或 `genpb.exe`）。

参数顺序：`[flag] [proto_in] [go_out] [cs_out] [lang] [tools_dir]`

| 参数 | 默认值 | 说明 |
|------|--------|------|
| flag | server | server \| client（client 不导出 data_srv/data_fwd） |
| proto_in | ./proto | 输入 proto 目录 |
| go_out | ./pb | Go 输出目录 |
| cs_out | ./pb/Pb | C# 输出目录 |
| lang | all | go \| Pb \| all |
| tools_dir | ../proto | 存放 protoc、protoc-gen-go 的目录 |

```bash
# Windows
genpb.bat                    # 全部默认
genpb.bat client             # client 模式
genpb.bat server ./proto ./pb ./pb/Pb go ../proto   # 指定工具目录

# Linux/macOS
./gen.sh
./gen.sh client
./gen.sh server ./proto ./pb ./pb/Pb go /opt/proto-tools
```

### Direct run

```bash
# Generate all (Go + C#)
go run main.go --proto_in ./proto

# Generate Go only
go run main.go --lang go --proto_in ./proto

# Generate C# only
go run main.go --lang Pb --proto_in ./proto
```

### Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `--lang` | `all` | Language: `go`, `Pb`, `all` |
| `--proto_in` | `./proto` | Proto input directory |
| `--go_out` | `./pb` | Go output directory |
| `--cs_out` | `./pb/Pb` | C# output directory |
| `--flag` | `server` | Export flag: `server` (all files), `client` (exclude `data_srv.proto`, `data_fwd.proto`) |
| `--tools_dir` | `../proto` | 存放 protoc、protoc-gen-go 的目录（空则用默认） |

## File Structure

```
genpb/
├── main.go      # Entry point + config
├── gen_go.go    # Go code generation
├── gen_cs.go    # C# code generation
├── gen.bat      # Windows 生成脚本
├── gen.sh       # Linux/macOS 生成脚本
├── CLAUDE.md    # This file
├── proto/       # Proto 定义
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
| `cmd_dsp.proto` | `Dsp*` | `DspLoginFast` → `EKey_LoginFast` |

### Dispatch (Dsp) Keys

In `cmd.proto`, keys defined between `// dsp start` and `// dsp end` comments are dispatch keys:

```protobuf
// dsp start
LoginFast          = 40000; // 同步快速重登 token
LoginData          = 40001; // 同步玩家登录数据
// dsp end
```

These keys map to `Dsp*` prefixed messages in `cmd_dsp.proto`.

## Related Projects

- `E:/tools/proto` - Protobuf source definitions
- `E:/tools/server` - Game server using generated types
- `E:/tools/client/sgame` - C# client
