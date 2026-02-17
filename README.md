# genpb

Protocol Buffer 代码生成工具。在 protoc 生成的基础代码之上，额外生成消息解析扩展与 Vector 定点数数学扩展，支持 **Go** 与 **C#**。

## 依赖

- **protoc**：Protocol Buffer 编译器
- **protoc-gen-go**：Go 语言插件（仅生成 Go 时需要）

默认从 `../proto` 目录查找上述可执行文件，可通过 `--tools_dir` 指定。

## 用法

```bash
go run .
```

或指定参数：

```bash
go run . --proto_in ./proto --go_out ./pb --cs_out ./pb/Pb --lang all --flag server
```

### 参数说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `--proto_in` | `./proto` | proto 文件所在目录 |
| `--go_out` | `./pb` | Go 生成输出目录 |
| `--cs_out` | `./pb/Pb` | C# 生成输出目录 |
| `--lang` | `all` | 生成语言：`go`、`Pb`、`all` |
| `--tools_dir` | `../proto` | 存放 protoc、protoc-gen-go 的目录 |
| `--flag` | `server` | 导出范围：`server`（全部）、`client`（排除 data_srv.proto、data_fwd.proto） |

## 生成内容

### Go

1. **protoc 生成**：`*.pb.go`（enum、data、cmd、cmd_req、cmd_rsp、cmd_dsp、data_srv、data_fwd 等）
2. **pbgen 扩展**：
   - **cmd.ext.go**：消息解析器（按 EKey 注册/反序列化）、`Marshal` / `Unmarshal`、各消息的 `Key()` / `Marshal()`
   - **data.pb.vector.go**：Vector 定点数数学扩展（见下方 Vector 定点数）

### C#

1. **protoc 生成**：`*.cs`（Data.cs、Enum.cs、Cmd.cs、CmdReq.cs、CmdRsp.cs、CmdDsp.cs 等）
2. **pbgen 扩展**：
   - **CmdExt.cs**：`Cmd.EKey` 枚举、`GetMessageType(EKey)`、各消息的 `GetKey()` 扩展方法
   - **DataVector.cs**：Vector 定点数数学扩展（与 Go 功能一致，见下方 Vector 定点数）

## Vector 定点数

`data.proto` 中定义：

```protobuf
message Vector {
  int64 x = 1;
  int64 y = 2;
  int64 z = 3;
}
```

Go 与 C# 的扩展均采用**定点数**表示坐标：

- **Scale = 1000**：真实坐标 `0.001` 对应存储值 `1`
- **常量**：`ZeroVector`、`ForwardVector`（面朝 X）、`OneVector`
- **构造/转换**：`NewVector(x,y,z)`（浮点→定点）、`NewVectorInt(x,y,z)`（定点）、`FloatToFixed`/`FixedToFloat`、`ToFloat64`/`Xf`/`Yf`/`Zf`
- **运算**：加减乘除（2D/3D、整数/浮点倍率）、点积/叉积、长度/距离、归一化、旋转、Lerp、MoveTowards、Min/Max/Clamp 等

Go 与 C# 的 Vector API 一一对应，便于服务端与客户端共用同一套坐标语义。

## 脚本

- **gen.bat**：Windows 下执行生成
- **gen.sh**：Linux/Mac 下执行生成
