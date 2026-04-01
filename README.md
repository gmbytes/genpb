# genpb

Protocol Buffer 代码生成工具。在 protoc 生成的基础代码之上，额外生成消息解析扩展与 Vector 定点数数学扩展，支持 **Go**、**C#** 与 **Rust**。

## 依赖

- **protoc**：Protocol Buffer 编译器
- **protoc-gen-go**：Go 语言插件（仅生成 Go 时需要）

默认从 `../proto` 目录查找上述可执行文件，可通过 `--tools_dir` 指定。

## 快速开始

先编译可执行文件：

```bash
go build -buildvcs=false -o genpb        # Linux/macOS
go build -buildvcs=false -o genpb.exe    # Windows
```

然后通过脚本或直接运行：

```bash
# Windows
gen.bat

# Linux/macOS
./gen.sh
```

或使用 `go run`：

```bash
go run .
```

## 用法

```bash
go run . --proto_in ./proto --go_out ./pb --cs_out ./pb/Pb --lang all --flag server
```

生成 Rust 客户端协议代码：

```bash
go run . --lang rust --flag client --rust_out ../../gclient/rust/netcore/src/gen
```

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--proto_in` | `./proto` | proto 文件所在目录 |
| `--go_out` | `./pb` | Go 生成输出目录 |
| `--cs_out` | `./pb/Pb` | C# 生成输出目录 |
| `--rust_out` | _(空)_ | Rust 生成输出目录（`--lang rust` 时必填，否则跳过 Rust 生成） |
| `--lang` | `all` | 生成语言：`go`、`Pb`、`rust`、`all` |
| `--tools_dir` | `../proto` | 存放 protoc、protoc-gen-go 的目录 |
| `--flag` | `server` | 导出范围：`server`（全部）、`client`（排除 data_srv.proto、data_fwd.proto） |

> 注意：`--lang rust` 时若未指定 `--rust_out` 会报错退出。

## 脚本

- **gen.bat**：Windows 下执行生成（调用编译好的 `genpb.exe`）
- **gen.sh**：Linux/macOS 下执行生成（调用编译好的 `genpb`）

脚本参数顺序：`[flag] [proto_in] [go_out] [cs_out] [lang] [tools_dir] [rust_out]`

| 位置 | 参数 | 默认值 | 说明 |
|------|------|--------|------|
| 1 | `flag` | `server` | `server` \| `client` |
| 2 | `proto_in` | `./proto` | 输入 proto 目录 |
| 3 | `go_out` | `./pb` | Go 输出目录（bat 默认指向 server/pb） |
| 4 | `cs_out` | `./pb/Pb` | C# 输出目录（bat 默认指向 client src） |
| 5 | `lang` | `all` | `go` \| `Pb` \| `rust` \| `all` |
| 6 | `tools_dir` | `../proto` | 存放 protoc 的目录 |
| 7 | `rust_out` | _(空)_ | Rust 输出目录（不设则跳过） |

示例：

```bash
# Windows
gen.bat                                             # 全部默认
gen.bat client                                      # client 模式
gen.bat server ./proto ./pb ./pb/Pb go ../proto     # 指定工具目录

# Linux/macOS
./gen.sh
./gen.sh client
./gen.sh server ./proto ./pb ./pb/Pb rust /opt/proto-tools ../../gclient/src/gen
```

## 文件结构

```
genpb/
├── main.go       # 入口 + 配置
├── gen_go.go     # Go 代码生成
├── gen_cs.go     # C# 代码生成
├── gen_rust.go   # Rust 代码生成
├── gen.bat       # Windows 生成脚本
├── gen.sh        # Linux/macOS 生成脚本
├── CLAUDE.md     # AI 辅助说明
├── proto/        # Proto 定义
└── pb/           # 生成文件（Go/C#）
```

## 生成内容

### Go

1. **protoc 生成**：`*.pb.go`（enum、data、cmd、cmd_req、cmd_rsp、cmd_dsp、data_srv、data_fwd 等）
2. **pbgen 扩展**：
   - **cmd.ext.go**：
     - `Package` 结构体：封装消息 + 错误码 + 缓存字节，二进制格式为 `[cmd 2B][errCode 2B][bodyLen 4B][body NB]`
     - `NewPackage(msg, errs...)` / `Key()` / `Marshal()`
     - `Unmarshal(key, data)` 全局反序列化入口
     - `parser`：按 EKey 注册/反序列化消息，全局 `_parser` 在 `init()` 中自动加载
     - 各消息的 `Key() EKey_T` 与 `Marshal() ([]byte, error)` 方法
   - **data.pb.vector.go**：Vector 定点数数学扩展（见下方 Vector 定点数）

### C#

1. **protoc 生成**：`*.cs`（Data.cs、Enum.cs、Cmd.cs、CmdReq.cs、CmdRsp.cs、CmdDsp.cs 等）
2. **pbgen 扩展**：
   - **CmdExt.cs**：`Cmd.EKey` 枚举、`GetMessageType(EKey)`、各消息的 `GetKey()` 扩展方法
   - **DataVector.cs**：Vector 定点数数学扩展（与 Go 功能一致，见下方 Vector 定点数）

### Rust

1. **cmd_ext.rs**：`EKey` 枚举（`#[repr(u16)]`）及其 `from_u16` / `as_u16`；`ClientMessage` / `ServerMessage` 枚举；`encode_client_message()` / `decode_server_message()`
2. **protocol.desc**：二进制 `FileDescriptorSet`，包含所有 proto message 定义，供运行时 descriptor 动态解码使用
3. **protocol_meta.json**：EKey 数值 → `{ ekey, message, event_name }` 的 JSON 映射表，供运行时通用通道查找消息类型

## Vector 定点数

`data.proto` 中定义：

```protobuf
message Vector {
  int64 x = 1;
  int64 y = 2;
  int64 z = 3;
}
```

Go 与 C# 的扩展均采用**定点数**表示坐标（X 轴为角色面朝方向，左手坐标系）：

- **Scale = 1000**：真实坐标 `0.001` 对应存储值 `1`
- **常量**：`ZeroVector`、`ForwardVector`（面朝 X 轴）、`OneVector`

### API 分类

| 类别 | 方法 |
|------|------|
| **构造/转换** | `NewVector(x,y,z)`（浮点→定点）、`NewVectorInt(x,y,z)`（定点）、`FloatToFixed` / `FixedToFloat`、`ToFloat64` / `Xf` / `Yf` / `Zf` |
| **字符串** | `StringF()`（真实浮点坐标字符串，格式 `(x.xxx, y.xxx, z.xxx)`） |
| **拷贝/反转** | `Copy()`、`CopyNewZ(z)`（替换 Z，浮点入参）、`CopyNewZInt(z)`（替换 Z，定点入参）、`CopyTo(dst)`、`Reverse2D()`、`Reverse()` |
| **加减** | `Add2D` / `Add`、`Sub2D` / `Sub`（2D 版本 Z 保持自身值） |
| **乘除** | `Mul2D` / `Mul`（整数倍率）、`MulFloat2D` / `MulFloat`（浮点倍率）、`Div2D` / `Div`（整数除，除零返回拷贝）、`DivFloat2D` / `DivFloat` |
| **点积/叉积** | `Dot2D` / `Dot`（定点数结果 = 真实值 × Scale²）、`Dot2DFloat` / `DotFloat`（真实浮点值）、`Cross`（三维叉积，自动 ÷Scale） |
| **长度** | `LengthSq2D` / `LengthSq`（定点数域平方）、`LengthSq2DFloat` / `LengthSqFloat`（真实浮点平方）、`Length2D` / `Length`（真实浮点长度） |
| **距离** | `DistanceSq2D` / `DistanceSq`、`DistanceSq2DFloat` / `DistanceSqFloat`、`Distance2D` / `Distance` |
| **比较** | `Equal2D` / `Equal`（精确比较）、`ApproximatelyEqual2D` / `ApproximatelyEqual`（容差 1 个定点单位，即真实 0.001） |
| **归一化** | `Norm2D()`（XOY 单位向量，保留 Z）、`Norm()`（三维单位向量） |
| **正交** | `Orthogonal2D()`（XOY 平面逆时针旋转 90°） |
| **角度/弧度** | `ToAngle2D()` / `Angle2D(v)`、`ToRadian2D()` / `Radian2D(v)` |
| **旋转** | `Rotate2D(rad)`（弧度，绕 Z 轴）、`RotateAngle2D(deg)`（角度） |
| **插值** | `Lerp(v, t)` / `Lerp2D(v, t)`（t ∈ [0,1]） |
| **移动** | `MoveTowards(target, maxDist)` / `MoveTowards2D(target, maxDist)`（maxDist 为真实浮点距离） |
| **随机** | `GenerateRandomVector(min, max)`（定点数范围内均匀随机） |
| **便捷** | `IsZero()` / `IsZero2D()`、`SetFromFloat64(x,y,z)`（就地修改）、`Min(v)` / `Max(v)` / `Clamp(min,max)`（分量操作） |

Go 与 C# 的 Vector API 一一对应，便于服务端与客户端共用同一套坐标语义。
