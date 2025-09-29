# Building lrcget-go for Multiple Platforms

This document explains how to build pre-compiled binaries for Windows, Linux, and macOS.

## Prerequisites

1. **Go 1.24+** - [Download here](https://golang.org/dl/)
2. **Node.js 20+** - [Download here](https://nodejs.org/)
3. **Wails CLI** - Install with: `go install github.com/wailsapp/wails/v2/cmd/wails@latest`

## Quick Build

### Using the build script (recommended)

**On macOS/Linux:**
```bash
./build.sh
```

**On Windows:**
```cmd
build.bat
```

### Manual build commands

```bash
# Build for current platform
wails build

# Build for specific platforms
wails build -platform windows/amd64 -o lrcget-windows-amd64.exe
wails build -platform linux/amd64 -o lrcget-linux-amd64
wails build -platform darwin/amd64 -o lrcget-darwin-amd64
wails build -platform darwin/arm64 -o lrcget-darwin-arm64
```

## Cross-Compilation Notes

### Windows to Linux/macOS
- Requires additional setup for cross-compilation
- Consider using GitHub Actions for reliable cross-platform builds

### macOS to Windows/Linux
- Requires cross-compilation toolchain setup
- GitHub Actions recommended for consistent builds

### Linux to Windows/macOS
- Requires cross-compilation setup
- GitHub Actions provides the most reliable results

## Automated Builds with GitHub Actions

The repository includes a GitHub Actions workflow (`.github/workflows/build.yml`) that automatically builds for all platforms when you create a release tag.

### Creating a release:

1. Create and push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. The workflow will automatically:
   - Build for all platforms
   - Create a GitHub release
   - Attach all binaries to the release

## Build Output

All binaries will be created in the `build/binaries/` directory:

- `lrcget-windows-amd64.exe` - Windows 64-bit
- `lrcget-linux-amd64` - Linux 64-bit
- `lrcget-darwin-amd64` - macOS Intel
- `lrcget-darwin-arm64` - macOS Apple Silicon

## Troubleshooting

### Common Issues:

1. **Wails not found**: Make sure Wails CLI is installed and in your PATH
2. **Node.js issues**: Ensure Node.js 20+ is installed
3. **Go version**: Ensure Go 1.24+ is installed
4. **Cross-compilation**: For reliable cross-platform builds, use GitHub Actions

### Platform-specific requirements:

- **Windows**: No additional requirements
- **Linux**: May need additional libraries for GUI support
- **macOS**: May need Xcode command line tools

## Distribution

The generated binaries are standalone executables that don't require users to install Go, Node.js, or any other dependencies. Users can simply download and run the appropriate binary for their platform.
