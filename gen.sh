#!/usr/bin/env bash
# 参数: [flag] [proto_in] [go_out] [cs_out] [lang] [tools_dir] [rust_out]
#   flag: server | client (默认 server)
#   proto_in: 输入 proto 目录 (默认 ./proto)
#   go_out: Go 输出目录 (默认 ./pb)
#   cs_out: C# 输出目录 (默认 ./pb/Pb)
#   lang: go | Pb | rust | all (默认 all)
#   tools_dir: 存放 protoc 和 protoc-gen-go 的目录 (默认 ../proto)
#   rust_out: Rust 输出目录 (默认空, 不设则跳过)

set -e
cd "$(dirname "$0")"

FLAG="${1:-server}"
PROTO_IN="${2:-./proto}"
GO_OUT="${3:-./pb}"
CS_OUT="${4:-./pb/Pb}"
LANG="${5:-all}"
TOOLS_DIR="${6:-../proto}"
RUST_OUT="${7:-}"

GENPB="./genpb"
if [[ ! -x "$GENPB" ]]; then
  GENPB="genpb"
fi

RUST_ARG=""
if [[ -n "$RUST_OUT" ]]; then
  RUST_ARG="--rust_out $RUST_OUT"
fi

echo "Generating protocol: flag=$FLAG, proto_in=$PROTO_IN, go_out=$GO_OUT, cs_out=$CS_OUT, lang=$LANG, tools=$TOOLS_DIR, rust_out=$RUST_OUT"
"$GENPB" --proto_in "$PROTO_IN" --go_out "$GO_OUT" --cs_out "$CS_OUT" --lang "$LANG" --flag "$FLAG" --tools_dir "$TOOLS_DIR" $RUST_ARG
