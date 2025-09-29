# LRCGET Go

A desktop utility for mass-downloading LRC synced lyrics for your music library, built with Go and Wails v2.

## Overview

LRCGET Go is a complete rewrite of the original LRCGET application, migrating from Rust/Tauri to Go/Wails v2. This provides better cross-platform compatibility and easier deployment while maintaining all the functionality of the original application.

## Features

- **Music Library Management**: Scan and organize your music collection
- **Lyrics Download**: Mass-download synced lyrics from LRCLIB
- **Audio Playback**: Built-in audio player with controls
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Modern UI**: Clean, responsive interface built with modern web technologies

## Architecture

### Backend (Go)
- **Database**: SQLite with migrations using `modernc.org/sqlite`
- **Audio**: Audio playback using `github.com/faiface/beep` and `github.com/hajimehoshi/oto`
- **HTTP Client**: LRCLIB API integration using `github.com/go-resty/resty`
- **File System**: Audio file scanning and metadata extraction using `github.com/dhowden/tag`

### Frontend (Web)
- **Framework**: Vanilla JavaScript with modern ES6+ features
- **Styling**: Custom CSS with dark/light theme support
- **Build**: Vite for fast development and optimized builds

## Project Structure

```
lrcget-go/
├── internal/                 # Internal Go packages
│   ├── app/                 # Main application logic
│   ├── audio/               # Audio player implementation
│   ├── database/            # Database layer with migrations
│   ├── filesystem/          # File system scanning
│   ├── lrclib/              # LRCLIB API client
│   └── utils/               # Utility functions
├── frontend/                # Frontend web application
│   ├── src/                 # Source files
│   └── dist/                # Built frontend
├── build/                   # Build assets and configurations
├── main.go                  # Application entry point
├── wails.json              # Wails configuration
└── go.mod                   # Go module dependencies
```

## System Requirements

### Minimum Requirements
- **OS**: Windows 10+, macOS 10.15+, or Linux (Ubuntu 18.04+)
- **RAM**: 4GB minimum, 8GB recommended
- **Storage**: 100MB for application, additional space for music library
- **Network**: Internet connection for lyrics download

### Development Requirements
- **Go**: 1.21 or later
- **Node.js**: 16.0 or later
- **npm**: 7.0 or later (included with Node.js)
- **Git**: 2.0 or later
- **Build Tools**: Platform-specific (see Installation section)

## Dependencies

### Core Dependencies
- **Wails v2**: Desktop application framework
- **SQLite**: Database with `modernc.org/sqlite`
- **Audio**: `github.com/faiface/beep` and `github.com/hajimehoshi/oto`
- **HTTP**: `github.com/go-resty/resty/v2`
- **Metadata**: `github.com/dhowden/tag`

### Development Dependencies
- **Vite**: Frontend build tool
- **Go 1.21+**: Go runtime
- **Node.js 16+**: JavaScript runtime
- **npm**: Node package manager

## Installation

### Prerequisites

#### Required Software
- **Go 1.21+**: Go programming language
- **Node.js 16+**: JavaScript runtime for frontend development
- **npm**: Node package manager (included with Node.js)
- **Wails v2 CLI**: Desktop application framework

#### Platform-Specific Requirements

##### Windows
- **Go**: Download from [golang.org](https://golang.org/dl/) or use [Chocolatey](https://chocolatey.org/): `choco install golang`
- **Node.js**: Download from [nodejs.org](https://nodejs.org/) or use Chocolatey: `choco install nodejs`
- **Git**: Download from [git-scm.com](https://git-scm.com/) or use Chocolatey: `choco install git`
- **Build Tools**: Visual Studio Build Tools or Visual Studio Community

##### macOS
- **Go**: Download from [golang.org](https://golang.org/dl/) or use [Homebrew](https://brew.sh/): `brew install go`
- **Node.js**: Download from [nodejs.org](https://nodejs.org/) or use Homebrew: `brew install node`
- **Git**: Usually pre-installed, or install via Homebrew: `brew install git`
- **Xcode Command Line Tools**: `xcode-select --install`

##### Linux (Ubuntu/Debian)
```bash
# Update package list
sudo apt update

# Install Go
sudo apt install golang-go

# Install Node.js (using NodeSource repository for latest version)
curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
sudo apt install nodejs

# Install Git
sudo apt install git

# Install build essentials
sudo apt install build-essential
```

##### Linux (CentOS/RHEL/Fedora)
```bash
# Install Go
sudo dnf install golang  # or yum install golang

# Install Node.js
curl -fsSL https://rpm.nodesource.com/setup_lts.x | sudo bash -
sudo dnf install nodejs  # or yum install nodejs

# Install Git
sudo dnf install git  # or yum install git

# Install build tools
sudo dnf groupinstall "Development Tools"  # or yum groupinstall "Development Tools"
```

### Setup Instructions

#### 1. Install Wails v2 CLI
```bash
# Install Wails v2 (works on all platforms)
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Add Go bin directory to PATH (if not already added)
# Windows: Add %USERPROFILE%\go\bin to PATH
# macOS/Linux: Add $HOME/go/bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin  # Add to ~/.bashrc, ~/.zshrc, etc.
```

#### 2. Clone the Repository
```bash
git clone https://github.com/your-username/lrcget-go.git
cd lrcget-go
```

#### 3. Install Dependencies

##### Go Dependencies
```bash
# Install Go module dependencies
go mod tidy
```

##### Frontend Dependencies
```bash
# Install Node.js dependencies
cd frontend
npm install
cd ..
```

#### 4. Verify Installation
```bash
# Check Go version
go version

# Check Node.js version
node --version
npm --version

# Check Wails version
wails version
```

### Platform-Specific Notes

#### Windows
- Ensure you have Visual Studio Build Tools installed
- You may need to restart your terminal after installing Go
- If you encounter permission issues, run your terminal as Administrator

#### macOS
- If you encounter issues with Xcode Command Line Tools, run: `sudo xcode-select --reset`
- You may need to accept Xcode license: `sudo xcodebuild -license accept`

#### Linux
- Some distributions may require additional packages for audio support
- For Ubuntu/Debian: `sudo apt install libasound2-dev`
- For CentOS/RHEL: `sudo dnf install alsa-lib-devel`

## Development

### Running in Development Mode
```bash
# Start development server with hot reload
wails dev

# The application will open automatically
# Frontend runs on http://localhost:5173/
# Backend API is available for frontend integration
```

### Building the Application

#### Development Build
```bash
# Build for current platform (development)
wails build -dev
```

#### Production Build
```bash
# Build for current platform (production)
wails build

# Build for specific platforms
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform darwin/arm64
wails build -platform linux/amd64
wails build -platform linux/arm64
```

#### Cross-Platform Building
```bash
# Build for multiple platforms at once
wails build -platform windows/amd64,darwin/amd64,linux/amd64

# Build with specific output directory
wails build -o ./dist
```

### Testing
```bash
# Run Go tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/database
go test ./internal/lrclib
go test ./internal/audio
```

### Troubleshooting

#### Common Issues

##### "Wails applications will not build without the correct build tags"
- **Solution**: Ensure you have the frontend built and Wails bindings generated
- **Fix**: Run `wails generate module` then `wails build`

##### "Could not resolve '../wailsjs/go/app/App'"
- **Solution**: Generate Wails bindings first
- **Fix**: Run `wails generate module` in the project root

##### Database migration errors
- **Solution**: Delete the existing database file and let the application recreate it
- **Fix**: Remove `~/.lrcget/db.sqlite3` and restart the application

##### Node.js/npm not found
- **Solution**: Ensure Node.js is properly installed and in PATH
- **Fix**: 
  - Windows: Restart terminal after installation
  - macOS/Linux: Check PATH with `echo $PATH` and add Node.js bin directory

##### Go module issues
- **Solution**: Ensure Go modules are properly initialized
- **Fix**: Run `go mod tidy` and `go mod download`

##### Audio playback issues (Linux)
- **Solution**: Install audio development libraries
- **Fix**: 
  - Ubuntu/Debian: `sudo apt install libasound2-dev`
  - CentOS/RHEL: `sudo dnf install alsa-lib-devel`

#### Platform-Specific Issues

##### Windows
- **Issue**: Build tools not found
- **Solution**: Install Visual Studio Build Tools or Visual Studio Community
- **Fix**: Download from Microsoft's website or use Chocolatey: `choco install visualstudio2019buildtools`

##### macOS
- **Issue**: Xcode Command Line Tools issues
- **Solution**: Reset and reinstall Xcode Command Line Tools
- **Fix**: `sudo xcode-select --reset && xcode-select --install`

##### Linux
- **Issue**: Missing development libraries
- **Solution**: Install build essentials and audio libraries
- **Fix**: 
  - Ubuntu/Debian: `sudo apt install build-essential libasound2-dev`
  - CentOS/RHEL: `sudo dnf groupinstall "Development Tools" && sudo dnf install alsa-lib-devel`

## Quick Start

### 1. Clone and Setup
```bash
# Clone the repository
git clone https://github.com/your-username/lrcget-go.git
cd lrcget-go

# Install dependencies
go mod tidy
cd frontend && npm install && cd ..

# Generate Wails bindings
wails generate module
```

### 2. Run Development Server
```bash
# Start development mode
wails dev
```

### 3. Build Application
```bash
# Build for your platform
wails build

# Run the built application
./build/bin/lrcget-go.app/Contents/MacOS/lrcget  # macOS
# or
./build/bin/lrcget-go.exe  # Windows
# or
./build/bin/lrcget-go  # Linux
```

## Usage

### First Time Setup
1. **Launch Application**: Start the application
2. **Configure Directories**: Add your music library directories
3. **Initialize Library**: Scan your music collection
4. **Download Lyrics**: Mass-download synced lyrics for your tracks

### Daily Usage
1. **Browse Library**: View your music by tracks, albums, or artists
2. **Play Music**: Use the built-in audio player
3. **Download Lyrics**: Get synced lyrics for tracks without them
4. **Manage Collection**: Organize and search your music library

### Features
- **Music Library Management**: Scan and organize your music collection
- **Lyrics Download**: Mass-download synced lyrics from LRCLIB
- **Audio Playback**: Built-in audio player with controls
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **Modern UI**: Clean, responsive interface

## Migration from Rust Version

This Go version maintains full compatibility with the original Rust version's database schema and functionality. The migration includes:

- **Database Schema**: Identical SQLite schema with all migrations
- **API Compatibility**: Same LRCLIB API integration
- **File Support**: Same audio format support
- **Features**: All original features preserved

## CI/CD and Deployment

### GitHub Actions (Recommended)
```yaml
# .github/workflows/build.yml
name: Build and Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v3
      with:
        go-version: '1.21'
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '18'
        cache: 'npm'
        cache-dependency-path: frontend/package-lock.json
    
    - name: Install dependencies
      run: |
        go mod tidy
        cd frontend && npm install && cd ..
    
    - name: Generate Wails bindings
      run: wails generate module
    
    - name: Build application
      run: wails build
```

### Docker Support
```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

# Install Node.js
RUN apk add --no-cache nodejs npm

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Install frontend dependencies
WORKDIR /app/frontend
RUN npm install

# Generate Wails bindings and build
WORKDIR /app
RUN wails generate module
RUN wails build

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/build/bin/lrcget-go .
CMD ["./lrcget-go"]
```

## Contributing

### Development Workflow
1. **Fork the repository**
2. **Create a feature branch**: `git checkout -b feature/your-feature-name`
3. **Make your changes**
4. **Add tests if applicable**
5. **Test on your platform**: `wails dev` and `wails build`
6. **Submit a pull request**

### Code Style
- **Go**: Follow standard Go formatting with `gofmt`
- **JavaScript**: Use modern ES6+ features
- **CSS**: Use consistent naming conventions
- **Commits**: Use conventional commit messages

### Testing Requirements
- **Unit Tests**: Add tests for new Go functions
- **Integration Tests**: Test API endpoints
- **Cross-Platform**: Test on multiple platforms if possible
- **Documentation**: Update README for new features

## License

This project is licensed under the same terms as the original LRCGET project.

## Acknowledgments

- Original LRCGET project for the concept and design
- LRCLIB for providing the lyrics API
- Wails team for the excellent desktop framework
- All contributors to the open-source libraries used