#!/usr/bin/env bash
# genpb 包装脚本（与 main.go 一致，已移除 C# / --cs_out）
#
# 参数: [flag] [proto_in] [go_out] [lang] [tools_dir] [rust_out] [godot_out]
#   flag      : server | client        (默认 server)
#   proto_in  : proto 目录             (默认 ./proto)
#   go_out    : Go 输出目录            (默认 ./pb)
#   lang      : go | rust | all        (默认 all)
#   tools_dir : protoc / protoc-gen-go / protoc-gen-prost (默认 ../proto)
#   rust_out  : Rust 输出目录；空则跳过 rust 生成
#   godot_out : godot_bridge_gen.rs；空则跳过
#
# 示例:
#   ./gen.sh client ./proto ./pb all ../proto \
#     ../../gclient/rust/lib/gnet/src/gen \
#     ../../gclient/rust/gdbridge/src/gen

set -e
cd "$(dirname "$0")"

FLAG="${1:-server}"
PROTO_IN="${2:-./proto}"
GO_OUT="${3:-./pb}"
LANG="${4:-all}"
TOOLS_DIR="${5:-../proto}"
RUST_OUT="${6:-}"
GODOT_OUT="${7:-}"

GENPB="./genpb"
if [[ ! -x "$GENPB" ]]; then
  GENPB="genpb"
fi

RUST_ARG=()
if [[ -n "$RUST_OUT" ]]; then
  RUST_ARG=(--rust_out "$RUST_OUT")
fi

GODOT_ARG=()
if [[ -n "$GODOT_OUT" ]]; then
  GODOT_ARG=(--godot_out "$GODOT_OUT")
fi

echo "Generating protocol: flag=$FLAG, proto_in=$PROTO_IN, go_out=$GO_OUT, lang=$LANG, tools=$TOOLS_DIR, rust_out=$RUST_OUT, godot_out=$GODOT_OUT"
"$GENPB" --proto_in "$PROTO_IN" --go_out "$GO_OUT" --lang "$LANG" --flag "$FLAG" --tools_dir "$TOOLS_DIR" "${RUST_ARG[@]}" "${GODOT_ARG[@]}"
