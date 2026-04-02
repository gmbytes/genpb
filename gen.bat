@echo off
chcp 65001 >nul 2>&1
setlocal
cd /d "%~dp0"

REM genpb wrapper (same as main.go; C# / --cs_out removed)
REM
REM Args: [flag] [proto_in] [go_out] [lang] [tools_dir] [rust_out] [godot_out]
REM   flag      : server | client        (default server)
REM   proto_in  : proto dir              (default ./proto)
REM   go_out    : Go output dir          (default ../../../server/server/internal/pb)
REM   lang      : go | rust | all        (default all)
REM   tools_dir : dir of protoc, protoc-gen-go, protoc-gen-prost (default .)
REM   rust_out  : Rust out; empty skips
REM   godot_out : godot_bridge_gen.rs out; empty skips; only with rust gen
REM
REM Examples:
REM   gen.bat
REM   gen.bat client
REM   gen.bat server ./proto ../../../server/server/internal/pb all . ..\..\..\gclient\rust\lib\gnet\src\gen ..\..\..\gclient\rust\gdbridge\src\gen

set FLAG=server
set PROTO_IN=./proto
set GO_OUT=../../../server/server/internal/pb
set LANG=all
set TOOLS_DIR=.
set RUST_OUT=
set GODOT_OUT=

if not "%~1"=="" set FLAG=%~1
if not "%~2"=="" set PROTO_IN=%~2
if not "%~3"=="" set GO_OUT=%~3
if not "%~4"=="" set LANG=%~4
if not "%~5"=="" set TOOLS_DIR=%~5
if not "%~6"=="" set RUST_OUT=%~6
if not "%~7"=="" set GODOT_OUT=%~7

set GENPB=%~dp0genpb.exe
if not exist "%GENPB%" set GENPB=genpb.exe

set RUST_ARG=
if not "%RUST_OUT%"=="" set RUST_ARG=--rust_out "%RUST_OUT%"

set GODOT_ARG=
if not "%GODOT_OUT%"=="" set GODOT_ARG=--godot_out "%GODOT_OUT%"

echo Generating protocol: flag=%FLAG%, proto_in=%PROTO_IN%, go_out=%GO_OUT%, lang=%LANG%, tools=%TOOLS_DIR%, rust_out=%RUST_OUT%, godot_out=%GODOT_OUT%
"%GENPB%" --proto_in "%PROTO_IN%" --go_out "%GO_OUT%" --lang %LANG% --flag %FLAG% --tools_dir "%TOOLS_DIR%" %RUST_ARG% %GODOT_ARG%

endlocal
