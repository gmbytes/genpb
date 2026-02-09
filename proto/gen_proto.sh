#!/bin/bash

# Protobuf 代码生成脚本
# 用法: ./gen_proto.sh

PROTO_DIR="proto"
OUT_DIR="proto/gen"
GO_OUT="$OUT_DIR/go"
CS_OUT="$OUT_DIR/cs"
PROTOC_VERSION="3.20.3"

# 工具路径
PROTOC="$PROTO_DIR/protoc"
PROTOC_GEN_GO="$PROTO_DIR/protoc-gen-go.exe"
PROTOC_GEN_GO_GRPC="$PROTO_DIR/protoc-gen-go-grpc.exe"
PROTOC_GEN_CS="$PROTO_DIR/protoc-gen-csharp.exe"

# 确保输出目录存在
mkdir -p "$GO_OUT"
mkdir -p "$CS_OUT"

echo "=== Protobuf 代码生成 ==="
echo "输入目录: $PROTO_DIR"
echo "Go 输出: $GO_OUT"
echo "C# 输出: $CS_OUT"

# 检查 protoc
if [ ! -f "$PROTOC" ]; then
    echo "错误: 未找到 protoc，请从 https://github.com/protocolbuffers/protobuf/releases 下载"
    exit 1
fi

# 检查并下载 protoc-gen-go
if [ ! -f "$PROTOC_GEN_GO" ]; then
    echo "下载 protoc-gen-go..."
    curl -L -o "$PROTO_DIR/protoc-gen-go.exe" \
        "https://github.com/protocolbuffers/protobuf-go/releases/download/v1.36.11/protoc-gen-go.exe"
fi

# 检查并下载 protoc-gen-go-grpc
if [ ! -f "$PROTOC_GEN_GO_GRPC" ]; then
    echo "下载 protoc-gen-go-grpc..."
    curl -L -o "$PROTO_DIR/protoc-gen-go-grpc.exe" \
        "https://github.com/grpc/grpc-go/releases/download/v1.60.0/protoc-gen-go-grpc.exe"
fi

# 检查并下载 protoc-gen-csharp
if [ ! -f "$PROTOC_GEN_CS" ]; then
    echo "下载 protoc-gen-csharp..."
    curl -L -o "$PROTO_DIR/protoc-gen-csharp.exe" \
        "https://github.com/protocolbuffers/protobuf-csharp/releases/download/v$PROTOC_VERSION/protoc-gen-csharp.exe"
fi

# 获取 proto 文件列表
PROTO_FILES=$(find "$PROTO_DIR" -maxdepth 1 -name "*.proto" -type f)

if [ -z "$PROTO_FILES" ]; then
    echo "错误: 未找到 .proto 文件"
    exit 1
fi

echo ""
echo "=== 生成 Go 代码 ==="
for proto in $PROTO_FILES; do
    echo "  处理: $(basename "$proto")"
    "$PROTOC" \
        --go_out="$GO_OUT" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$GO_OUT" \
        --go-grpc_opt=paths=source_relative \
        -I "$PROTO_DIR" \
        -I "$PROTO_DIR/include" \
        "$proto" 2>&1
    if [ $? -ne 0 ]; then
        echo "错误: Go 代码生成失败: $proto"
    fi
done

echo ""
echo "=== 生成 C# 代码 ==="
for proto in $PROTO_FILES; do
    echo "  处理: $(basename "$proto")"
    "$PROTOC" \
        --csharp_out="$CS_OUT" \
        --grpc_out="$CS_OUT" \
        --plugin=protoc-gen-grpc="$PROTOC_GEN_CS" \
        -I "$PROTO_DIR" \
        -I "$PROTO_DIR/include" \
        "$proto" 2>&1
    if [ $? -ne 0 ]; then
        echo "错误: C# 代码生成失败: $proto"
    fi
done

echo ""
echo "=== 完成 ==="
echo "Go 代码已生成到: $GO_OUT"
echo "C# 代码已生成到: $CS_OUT"
