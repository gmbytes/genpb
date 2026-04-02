# genpb

Protocol Buffer 代码生成工具。在 `protoc` 生成的基础代码之上，额外生成消息解析扩展，支持 **Go** 与 **Rust/Godot**。  
（C# 已废弃，相关代码已移除。）

## 架构概览

```
proto files
  │
  ├──► protoc (gen_go.go)         ──► Go pb / cmd.ext.go    (server 侧，保持不变)
  │
  └──► protoc (gen_rust.go)
         │  --descriptor_set_out  ──► protocol.desc          (热更 fallback 通道)
         │  --prost_out           ──► pb.rs                  (prost 类型)
         │  buildProtocolManifest ──► protocol_manifest.json (协议索引)
         │  writeTypedProtocol    ──► typed_protocol.rs      (EKey / ServerMessage / ClientMessage)
         └──► writeGodotBridgeGen ──► godot_bridge_gen.rs    (GodotClass + mapper，需 --godot_out)
```

- `genpb` 的 Rust 侧以 **`ProtocolManifest`** 为单一真源，来自 `FileDescriptorSet`（descriptor），取代了之前的正则解析。
- `protocol.desc` 保留，退为热更/fallback 旁路通道（不再是主数据源）。
- Go 侧 `gen_go.go` 完全不动。

## 依赖

- **protoc**：Protocol Buffer 编译器
- **protoc-gen-go**：Go 语言插件（仅 Go 生成需要）
- **protoc-gen-prost**：Rust 语言插件（仅 Rust 生成需要）
  ```bash
  cargo install protoc-gen-prost
  ```

默认从 `../proto` 目录查找上述可执行文件，可通过 `--tools_dir` 指定。

## 快速开始

```bash
go build -buildvcs=false -o genpb       # Linux/macOS
go build -buildvcs=false -o genpb.exe   # Windows
```

或直接用 `go run`：

```bash
# 生成 Go（服务端）
go run -buildvcs=false . --lang go --flag server --go_out ./pb

# 生成 Rust + Godot bridge（客户端）
go run -buildvcs=false . --lang rust --flag client \
    --rust_out  ../../gclient/rust/lib/gnet/src/gen \
    --godot_out ../../gclient/rust/gdbridge/src/gen
```

## 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--proto_in`  | `./proto` | proto 文件所在目录 |
| `--go_out`    | `./pb`    | Go 生成输出目录 |
| `--rust_out`  | _(空)_    | Rust 生成输出目录（`pb.rs`、`typed_protocol.rs`、`protocol_manifest.json`、`protocol.desc`） |
| `--godot_out` | _(空)_    | Godot bridge 输出目录（`godot_bridge_gen.rs`）；与 `--rust_out` 搭配使用 |
| `--lang`      | `all`     | 生成语言：`go`、`rust`、`all` |
| `--tools_dir` | `../proto`| 存放 `protoc`、`protoc-gen-go`、`protoc-gen-prost` 的目录（或确保它们在 `PATH` 中） |
| `--flag`      | `server`  | 导出范围：`server`（全部）、`client`（排除 data_srv.proto、data_fwd.proto） |

> `--lang rust` 时若未指定 `--rust_out` 会报错退出。  
> `--godot_out` 可选；不设则跳过 `godot_bridge_gen.rs` 生成。  
> **C# 已废弃**：不再支持 `--lang Pb`、`--cs_out`。

### 批处理 / Shell 脚本

本目录 `gen.bat`、`gen.sh` 已与 `main.go` 对齐（无 `--cs_out`，支持 `--rust_out` / `--godot_out`）。

参数顺序：`[flag] [proto_in] [go_out] [lang] [tools_dir] [rust_out] [godot_out]`（后两项可省略）。

Windows 示例（在 `comm/tools/genpb` 下，且已 `go build -o genpb.exe`）：

```bat
gen.bat server .\proto ..\..\..\server\server\internal\pb all . ..\..\..\gclient\rust\lib\gnet\src\gen ..\..\..\gclient\rust\gdbridge\src\gen
```

仓库根上级的 `comm/genpb.bat` 从 `comm/` 调用 `tools/genpb/genpb.exe`，参数顺序相同；默认只生成 Go，需要 Rust/Godot 时传入第 6、7 项路径。

等价 `go run` 示例：

```bash
go run -buildvcs=false . --lang rust --flag client \
  --rust_out  ../../gclient/rust/lib/gnet/src/gen \
  --godot_out ../../gclient/rust/gdbridge/src/gen
```

## 文件结构

```
genpb/
├── main.go       # 入口 + 配置
├── manifest.go   # ProtocolManifest 数据模型 + FileDescriptorSet 解析
├── gen_go.go     # Go 代码生成（服务端，不变）
├── gen_rust.go   # Rust/Godot 代码生成（manifest 驱动）
├── gen.bat       # Windows 生成脚本
├── gen.sh        # Linux/macOS 生成脚本
├── CLAUDE.md     # AI 辅助说明
├── proto/        # Proto 定义
└── pb/           # Go 生成文件
```

## 生成内容

### Go（服务端，不变）

1. **protoc 生成**：`*.pb.go`
2. **cmd.ext.go**：`Package` 结构体、`Unmarshal`、消息 `Key()` / `Marshal()` 扩展

### Rust（`--rust_out`）

| 文件 | 说明 |
|------|------|
| `pb.rs` | `protoc-gen-prost` 生成的 prost 类型，静态 `include!` 到 gnet |
| `typed_protocol.rs` | `EKey` 枚举（`#[repr(u16)]`）、`ServerMessage` / `ClientMessage` 枚举、`encode_client_message` / `decode_server_message`、`event_name()` 方法、`COMPILED_FINGERPRINTS`（递归 FNV-1a 指纹）|
| `protocol_manifest.json` | 协议索引（替代旧 `protocol_meta.json`）：每条 EKey 的方向、kind、event_name、字段 schema、fingerprint 等 |
| `protocol.desc` | 二进制 `FileDescriptorSet`，供热更 fallback 通道（`prost-reflect`）使用 |

### Godot bridge（`--godot_out`，可选）

| 文件 | 说明 |
|------|------|
| `godot_bridge_gen.rs` | 数据类 GodotClass（`VectorGd`、`TimeGd` 等）、服务端消息 GodotClass（`RspLoginGd`、`DspMoveGd` 等）、`NetEventGd` 统一事件包装、`server_message_to_event()` dispatch 函数、`hotfix_to_event()` fallback 函数 |

## Vector 定点数

`data.proto` 中定义 `Vector { int64 x, y, z }`，Go 侧有完整定点数扩展（Scale = 1000）。
详见 `gen_go.go` 生成的 `data.pb.vector.go`。
