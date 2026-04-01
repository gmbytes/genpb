@echo off
chcp 65001 >nul 2>&1
setlocal
cd /d "%~dp0"

REM Parameters: [flag] [proto_in] [go_out] [cs_out] [lang] [tools_dir] [rust_out]
REM   flag: server | client (default: server)
REM   proto_in: Input proto directory (default: ./proto)
REM   go_out: Go output directory (default: ./pb)
REM   cs_out: C# output directory (default: ./pb/Pb)
REM   lang: go | Pb | rust | all (default: all)
REM   tools_dir: Directory containing protoc.exe and protoc-gen-go.exe (default: ../proto)
REM   rust_out: Rust output directory (default: empty, skipped if not set)

set FLAG=server
set PROTO_IN=./proto
set GO_OUT=../../../server/server/internal/pb
set CS_OUT=../../../client/src/net/pb/
set LANG=all
set TOOLS_DIR=.
set RUST_OUT=

if not "%~1"=="" set FLAG=%~1
if not "%~2"=="" set PROTO_IN=%~2
if not "%~3"=="" set GO_OUT=%~3
if not "%~4"=="" set CS_OUT=%~4
if not "%~5"=="" set LANG=%~5
if not "%~6"=="" set TOOLS_DIR=%~6
if not "%~7"=="" set RUST_OUT=%~7

set GENPB=%~dp0genpb.exe
if not exist "%GENPB%" set GENPB=genpb.exe

set RUST_ARG=
if not "%RUST_OUT%"=="" set RUST_ARG=--rust_out "%RUST_OUT%"

echo Generating protocol: flag=%FLAG%, proto_in=%PROTO_IN%, go_out=%GO_OUT%, cs_out=%CS_OUT%, lang=%LANG%, tools=%TOOLS_DIR%, rust_out=%RUST_OUT%
"%GENPB%" --proto_in "%PROTO_IN%" --go_out "%GO_OUT%" --cs_out "%CS_OUT%" --lang %LANG% --flag %FLAG% --tools_dir "%TOOLS_DIR%" %RUST_ARG%

endlocal
