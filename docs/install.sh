#!/bin/bash
set -e

# Ollama MCP Server Installation Script
# Usage: curl -fsSL https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.sh | bash

REPO="efortin/ollama-mcp"
BINARY_NAME="ollama-mcp-server"
INSTALL_DIR="${OLLAMA_MCP_INSTALL_DIR:-$HOME/.local/bin}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
info() {
    printf "${GREEN}[INFO]${NC} %s\n" "$1" >&2
}

error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1" >&2
    exit 1
}

warn() {
    printf "${YELLOW}[WARN]${NC} %s\n" "$1" >&2
}

# Detect OS and architecture
detect_platform() {
    local os=""
    local arch=""
    
    # Detect OS
    case "$(uname -s)" in
        Linux*)     os="linux";;
        Darwin*)    os="darwin";;
        CYGWIN*|MINGW*|MSYS*) os="windows";;
        *)          error "Unsupported operating system: $(uname -s)";;
    esac
    
    # Detect architecture
    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64";;
        arm64|aarch64)  arch="arm64";;
        *)              error "Unsupported architecture: $(uname -m)";;
    esac
    
    echo "${os}-${arch}"
}

# Get the latest release version
get_latest_version() {
    local latest_release_url="https://api.github.com/repos/${REPO}/releases/latest"
    local version
    
    info "Fetching latest release information..."
    
    if command -v curl >/dev/null 2>&1; then
        # Use -f to fail on HTTP errors and --max-time for timeout
        version=$(curl -s -f --max-time 10 "$latest_release_url" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    elif command -v wget >/dev/null 2>&1; then
        version=$(wget -qO- --timeout=10 "$latest_release_url" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    else
        error "Neither curl nor wget found. Please install one of them."
        return 1
    fi
    
    # Check if we got a valid version
    if [ -z "$version" ]; then
        # Try to get the most recent release from the releases list
        info "No 'latest' release found, checking all releases..."
        local releases_url="https://api.github.com/repos/${REPO}/releases"
        
        if command -v curl >/dev/null 2>&1; then
            version=$(curl -s -f --max-time 10 "$releases_url" 2>/dev/null | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
        else
            version=$(wget -qO- --timeout=10 "$releases_url" 2>/dev/null | grep '"tag_name":' | head -1 | sed -E 's/.*"([^"]+)".*/\1/')
        fi
        
        if [ -z "$version" ]; then
            warn "No releases found on GitHub"
            return 1
        fi
    fi
    
    echo "$version"
    return 0
}

# Download the binary
download_binary() {
    local version="$1"
    local platform="$2"
    local binary_name="${BINARY_NAME}-${version#v}-${platform}"
    
    # Add .exe extension for Windows
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="${binary_name}.exe"
    fi
    
    local download_url="https://github.com/${REPO}/releases/download/${version}/${binary_name}"
    local temp_file="/tmp/${BINARY_NAME}-download"
    
    info "Downloading ${binary_name}..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$download_url" -o "$temp_file" || error "Failed to download binary"
    else
        wget -q "$download_url" -O "$temp_file" || error "Failed to download binary"
    fi
    
    echo "$temp_file"
}

# Verify checksum
verify_checksum() {
    local version="$1"
    local platform="$2"
    local binary_file="$3"
    local binary_name="${BINARY_NAME}-${version#v}-${platform}"
    
    # Add .exe extension for Windows
    if [[ "$platform" == *"windows"* ]]; then
        binary_name="${binary_name}.exe"
    fi
    
    info "Verifying checksum..."
    
    local checksums_url="https://github.com/${REPO}/releases/download/${version}/checksums.txt"
    local temp_checksums="/tmp/ollama-mcp-checksums.txt"
    
    if command -v curl >/dev/null 2>&1; then
        curl -fsSL "$checksums_url" -o "$temp_checksums" || warn "Failed to download checksums"
    else
        wget -q "$checksums_url" -O "$temp_checksums" || warn "Failed to download checksums"
    fi
    
    if [ -f "$temp_checksums" ]; then
        local expected_checksum=$(grep "$binary_name" "$temp_checksums" | awk '{print $1}')
        if [ -n "$expected_checksum" ]; then
            local actual_checksum=""
            if command -v sha256sum >/dev/null 2>&1; then
                actual_checksum=$(sha256sum "$binary_file" | awk '{print $1}')
            elif command -v shasum >/dev/null 2>&1; then
                actual_checksum=$(shasum -a 256 "$binary_file" | awk '{print $1}')
            else
                warn "No checksum tool found, skipping verification"
                rm -f "$temp_checksums"
                return
            fi
            
            if [ "$expected_checksum" != "$actual_checksum" ]; then
                rm -f "$temp_checksums"
                error "Checksum verification failed"
            fi
            info "Checksum verified successfully"
        else
            warn "Checksum not found for $binary_name"
        fi
        rm -f "$temp_checksums"
    fi
}

# Install the binary
install_binary() {
    local binary_file="$1"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    # Move binary to install directory
    local target_path="${INSTALL_DIR}/${BINARY_NAME}"
    mv "$binary_file" "$target_path"
    chmod +x "$target_path"
    
    info "Installed to: $target_path"
    
    # Check if install directory is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "$INSTALL_DIR is not in your PATH"
        echo ""
        echo "Add it to your PATH by adding this line to your shell configuration file:"
        echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
    fi
}

# Main installation process
main() {
    echo "Ollama MCP Server Installer"
    echo "=========================="
    echo ""
    
    # Check if running on Windows without proper shell
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" ]]; then
        warn "Windows detected. For best results, use PowerShell or WSL."
    fi
    
    # Detect platform
    local platform=$(detect_platform)
    info "Detected platform: $platform"
    
    # Get latest version or use specified version
    local version=""
    if [ -n "$OLLAMA_MCP_VERSION" ]; then
        version="$OLLAMA_MCP_VERSION"
        info "Using specified version: $version"
    else
        version=$(get_latest_version)
        if [ $? -ne 0 ] || [ -z "$version" ]; then
            error "No releases available yet. Please check https://github.com/${REPO}/releases"
            echo ""
            echo "Alternatively, you can build from source:"
            echo "  git clone https://github.com/${REPO}.git"
            echo "  cd ollama-mcp"
            echo "  mise run build"
            exit 1
        fi
    fi
    
    # Ensure version has 'v' prefix
    if [[ "$version" != v* ]]; then
        version="v$version"
    fi
    
    info "Installing version: $version"
    
    # Download binary
    local binary_file=$(download_binary "$version" "$platform")
    
    # Verify checksum
    verify_checksum "$version" "$platform" "$binary_file"
    
    # Install binary
    install_binary "$binary_file"
    
    # Test installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        info "Installation successful!"
        echo ""
        "$BINARY_NAME" --version
    else
        info "Installation complete!"
        echo ""
        echo "Run the following command to verify:"
        echo "  ${INSTALL_DIR}/${BINARY_NAME} --version"
    fi
    
    echo ""
    echo "To use with MCP, add the server to your MCP configuration."
    echo "For more information, visit: https://github.com/${REPO}"
}

# Run main function
main "$@"
