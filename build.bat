@echo off
REM Build script for lrcget-go for Windows
REM This script builds the application for Windows

echo 🚀 Building lrcget-go for Windows...

REM Check if wails is installed
where wails >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo ❌ Wails CLI not found. Please install it first:
    echo    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    exit /b 1
)

REM Create build directory
if not exist build\binaries mkdir build\binaries

REM Build for Windows (64-bit)
echo 🪟 Building for Windows (amd64)...
wails build -platform windows/amd64 -clean -o build\binaries\lrcget-windows-amd64.exe

REM Build for Linux (64-bit) - requires cross-compilation setup
echo 🐧 Building for Linux (amd64)...
wails build -platform linux/amd64 -clean -o build\binaries\lrcget-linux-amd64

REM Build for macOS (64-bit) - requires cross-compilation setup
echo 🍎 Building for macOS (amd64)...
wails build -platform darwin/amd64 -clean -o build\binaries\lrcget-darwin-amd64

REM Build for macOS (Apple Silicon) - requires cross-compilation setup
echo 🍎 Building for macOS (arm64)...
wails build -platform darwin/arm64 -clean -o build\binaries\lrcget-darwin-arm64

echo ✅ Build completed! Binaries are in build\binaries\
echo.
echo Generated binaries:
dir build\binaries\

pause
