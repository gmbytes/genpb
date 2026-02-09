@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM ============================================
REM Protobuf 代码生成脚本 (Windows)
REM 用法: gen_proto.bat
REM ============================================

set "PROTO_DIR=proto"
set "OUT_DIR=proto/gen"
set "GO_OUT=genpb/pb"
set "CS_OUT=%OUT_DIR%/cs"

REM 工具路径 - 使用 proto 目录下的工具
set "PROTOC=%PROTO_DIR%/protoc.exe"
set "PROTOC_GEN_GO=%PROTO_DIR%/protoc-gen-go.exe"

REM 创建输出目录
if not exist "%GO_OUT%" mkdir "%GO_OUT%"
if not exist "%CS_OUT%" mkdir "%CS_OUT%"

echo ===============================================
echo Protobuf 代码生成
echo ===============================================
echo 输入目录: %PROTO_DIR%
echo Go 输出: %GO_OUT%
echo C# 输出: %CS_OUT%

REM 检查 protoc
if not exist "%PROTOC%" (
    echo 错误: 未找到 protoc.exe
    exit /b 1
)

REM 检查 protoc-gen-go
if not exist "%PROTOC_GEN_GO%" (
    echo 错误: 未找到 protoc-gen-go.exe
    exit /b 1
)

echo.
echo ===============================================
echo 生成 Go 代码
echo ===============================================
for %%f in (%PROTO_DIR%\*.proto) do (
    echo   处理: %%~nf%%~xf
    "%PROTOC%" ^
        --go_out="%GO_OUT%" ^
        --go_opt=paths=source_relative ^
        -I "%PROTO_DIR%" ^
        "%%f"
    if errorlevel 1 (
        echo 警告: Go 代码生成失败: %%f
    )
)

echo.
echo ===============================================
echo 生成 C# 代码
echo ===============================================
for %%f in (%PROTO_DIR%\*.proto) do (
    set "SKIP=0"
    if "%%~nxf"=="data_fwd.proto" set SKIP=1
    if "%%~nxf"=="data_srv.proto" set SKIP=1
    if "!SKIP!"=="1" (
        echo   跳过: %%~nxf
    ) else (
        echo   处理: %%~nf%%~xf
        "%PROTOC%" ^
            --csharp_out="%CS_OUT%" ^
            -I "%PROTO_DIR%" ^
            "%%f"
        if errorlevel 1 (
            echo 警告: C# 代码生成失败: %%f
        )
    )
)

echo.
echo ===============================================
echo 完成
echo ===============================================
echo Go 代码已生成到: %GO_OUT%
echo C# 代码已生成到: %CS_OUT%
endlocal
pause
