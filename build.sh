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

# Check if frontend dependencies are installed
if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend
    npm install
    cd ..
fi

# Create build directory
mkdir -p build/binaries

# Function to build and verify
build_platform() {
    local platform=$1
    local output=$2
    local description=$3
    
    echo "🔨 Building for $description..."
    if wails build -platform $platform -clean -o build/binaries/$output; then
        if [ -f "build/binaries/$output" ]; then
            echo "✅ $description build successful"
            ls -lh build/binaries/$output
        else
            echo "❌ $description build failed - no output file"
            return 1
        fi
    else
        echo "❌ $description build failed"
        return 1
    fi
}

# Build for current platform first (to test setup)
echo "📦 Testing build for current platform..."
if wails build -clean; then
    echo "✅ Current platform build successful"
else
    echo "❌ Current platform build failed - check your setup"
    exit 1
fi

# Build for all platforms
echo ""
echo "🌍 Building for all platforms..."

# Build for Windows (64-bit)
build_platform "windows/amd64" "lrcget-windows-amd64.exe" "Windows (amd64)"

# Build for Linux (64-bit)
build_platform "linux/amd64" "lrcget-linux-amd64" "Linux (amd64)"

# Build for macOS (64-bit)
build_platform "darwin/amd64" "lrcget-darwin-amd64" "macOS Intel (amd64)"

# Build for macOS (Apple Silicon)
build_platform "darwin/arm64" "lrcget-darwin-arm64" "macOS Apple Silicon (arm64)"

echo ""
echo "✅ All builds completed! Binaries are in build/binaries/"
echo ""
echo "Generated binaries:"
ls -la build/binaries/
echo ""
echo "File sizes:"
ls -lh build/binaries/
