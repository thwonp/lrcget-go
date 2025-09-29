#!/bin/bash

# Build script for lrcget-go for multiple platforms
# This script builds the application for Windows, Linux, and macOS

set -e

echo "🚀 Building lrcget-go for multiple platforms..."

# Check if wails is installed
if ! command -v wails &> /dev/null; then
    echo "❌ Wails CLI not found. Please install it first:"
    echo "   go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# Create build directory
mkdir -p build/binaries

# Build for current platform (development)
echo "📦 Building for current platform..."
wails build -clean

# Build for Windows (64-bit)
echo "🪟 Building for Windows (amd64)..."
wails build -platform windows/amd64 -clean -o build/binaries/lrcget-windows-amd64.exe

# Build for Linux (64-bit)
echo "🐧 Building for Linux (amd64)..."
wails build -platform linux/amd64 -clean -o build/binaries/lrcget-linux-amd64

# Build for macOS (64-bit)
echo "🍎 Building for macOS (amd64)..."
wails build -platform darwin/amd64 -clean -o build/binaries/lrcget-darwin-amd64

# Build for macOS (Apple Silicon)
echo "🍎 Building for macOS (arm64)..."
wails build -platform darwin/arm64 -clean -o build/binaries/lrcget-darwin-arm64

echo "✅ Build completed! Binaries are in build/binaries/"
echo ""
echo "Generated binaries:"
ls -la build/binaries/
