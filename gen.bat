@echo off
setlocal
cd /d "%~dp0"

REM 参数: [flag] [proto_in] [go_out] [cs_out] [lang] [tools_dir]
REM  flag: server | client (默认 server)
REM  proto_in: 输入 proto 目录 (默认 ./proto)
REM  go_out: Go 输出目录 (默认 ./pb)
REM  cs_out: C# 输出目录 (默认 ./pb/Pb)
REM  lang: go | Pb | all (默认 all)
REM  tools_dir: 存放 protoc.exe 和 protoc-gen-go.exe 的目录 (默认 ../proto)

set FLAG=server
set PROTO_IN=./proto
set GO_OUT=./pb
set CS_OUT=./pb/Pb
set LANG=all
set TOOLS_DIR=../proto

if not "%~1"=="" set FLAG=%~1
if not "%~2"=="" set PROTO_IN=%~2
if not "%~3"=="" set GO_OUT=%~3
if not "%~4"=="" set CS_OUT=%~4
if not "%~5"=="" set LANG=%~5
if not "%~6"=="" set TOOLS_DIR=%~6

set GENPB=%~dp0genpb.exe
if not exist "%GENPB%" set GENPB=genpb.exe

echo Generating protocol: flag=%FLAG%, proto_in=%PROTO_IN%, go_out=%GO_OUT%, cs_out=%CS_OUT%, lang=%LANG%, tools=%TOOLS_DIR%
"%GENPB%" --proto_in "%PROTO_IN%" --go_out "%GO_OUT%" --cs_out "%CS_OUT%" --lang %LANG% --flag %FLAG% --tools_dir "%TOOLS_DIR%"

endlocal
