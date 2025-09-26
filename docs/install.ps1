# Ollama MCP Server Installation Script for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.ps1 | iex

$ErrorActionPreference = "Stop"

$REPO = "efortin/ollama-mcp"
$BINARY_NAME = "ollama-mcp-server"
$DEFAULT_INSTALL_DIR = "$env:LOCALAPPDATA\Programs\ollama-mcp"

# Use custom install directory if provided
$INSTALL_DIR = if ($env:OLLAMA_MCP_INSTALL_DIR) { 
    $env:OLLAMA_MCP_INSTALL_DIR 
} else { 
    $DEFAULT_INSTALL_DIR 
}

# Helper functions
function Write-Info {
    param($Message)
    Write-Host "[INFO] " -ForegroundColor Green -NoNewline
    Write-Host $Message
}

function Write-Error {
    param($Message)
    Write-Host "[ERROR] " -ForegroundColor Red -NoNewline
    Write-Host $Message
    exit 1
}

function Write-Warn {
    param($Message)
    Write-Host "[WARN] " -ForegroundColor Yellow -NoNewline
    Write-Host $Message
}

# Get the latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest release information..."
    
    try {
        $releases = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $version = $releases.tag_name
        
        if (-not $version) {
            Write-Error "Failed to fetch latest version"
        }
        
        return $version
    }
    catch {
        Write-Error "Failed to fetch release information: $_"
    }
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    
    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        default { Write-Error "Unsupported architecture: $arch" }
    }
}

# Download the binary
function Download-Binary {
    param(
        [string]$Version,
        [string]$Architecture
    )
    
    $binaryName = "$BINARY_NAME-$($Version.TrimStart('v'))-windows-$Architecture.exe"
    $downloadUrl = "https://github.com/$REPO/releases/download/$Version/$binaryName"
    $tempFile = Join-Path $env:TEMP "$BINARY_NAME-download.exe"
    
    Write-Info "Downloading $binaryName..."
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile -UseBasicParsing
        return $tempFile
    }
    catch {
        Write-Error "Failed to download binary: $_"
    }
}

# Verify checksum
function Verify-Checksum {
    param(
        [string]$Version,
        [string]$Architecture,
        [string]$BinaryFile
    )
    
    Write-Info "Verifying checksum..."
    
    $binaryName = "$BINARY_NAME-$($Version.TrimStart('v'))-windows-$Architecture.exe"
    $checksumsUrl = "https://github.com/$REPO/releases/download/$Version/checksums.txt"
    $tempChecksums = Join-Path $env:TEMP "ollama-mcp-checksums.txt"
    
    try {
        Invoke-WebRequest -Uri $checksumsUrl -OutFile $tempChecksums -UseBasicParsing
        
        $checksums = Get-Content $tempChecksums
        $expectedLine = $checksums | Where-Object { $_ -match $binaryName }
        
        if ($expectedLine) {
            $expectedChecksum = ($expectedLine -split '\s+')[0]
            $actualChecksum = (Get-FileHash -Path $BinaryFile -Algorithm SHA256).Hash.ToLower()
            
            if ($expectedChecksum -ne $actualChecksum) {
                Remove-Item $tempChecksums -Force
                Write-Error "Checksum verification failed"
            }
            
            Write-Info "Checksum verified successfully"
        }
        else {
            Write-Warn "Checksum not found for $binaryName"
        }
        
        Remove-Item $tempChecksums -Force
    }
    catch {
        Write-Warn "Failed to verify checksum: $_"
    }
}

# Install the binary
function Install-Binary {
    param(
        [string]$BinaryFile
    )
    
    # Create install directory if it doesn't exist
    if (-not (Test-Path $INSTALL_DIR)) {
        New-Item -ItemType Directory -Path $INSTALL_DIR -Force | Out-Null
    }
    
    # Move binary to install directory
    $targetPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    Move-Item -Path $BinaryFile -Destination $targetPath -Force
    
    Write-Info "Installed to: $targetPath"
    
    # Check if install directory is in PATH
    $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($userPath -notlike "*$INSTALL_DIR*") {
        Write-Warn "$INSTALL_DIR is not in your PATH"
        Write-Host ""
        Write-Host "Would you like to add it to your PATH? (Y/n): " -NoNewline
        $response = Read-Host
        
        if ($response -eq '' -or $response -eq 'Y' -or $response -eq 'y') {
            $newPath = $userPath + ";$INSTALL_DIR"
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
            Write-Info "Added to PATH. Please restart your terminal for changes to take effect."
        }
        else {
            Write-Host ""
            Write-Host "To add it manually, run:"
            Write-Host "  `$env:Path += `";$INSTALL_DIR`""
            Write-Host ""
        }
    }
}

# Main installation process
function Main {
    Write-Host "Ollama MCP Server Installer for Windows"
    Write-Host "======================================="
    Write-Host ""
    
    # Detect architecture
    $architecture = Get-Architecture
    Write-Info "Detected architecture: $architecture"
    
    # Get latest version or use specified version
    $version = if ($env:OLLAMA_MCP_VERSION) { 
        $env:OLLAMA_MCP_VERSION 
    } else { 
        Get-LatestVersion 
    }
    Write-Info "Installing version: $version"
    
    # Download binary
    $binaryFile = Download-Binary -Version $version -Architecture $architecture
    
    # Verify checksum
    Verify-Checksum -Version $version -Architecture $architecture -BinaryFile $binaryFile
    
    # Install binary
    Install-Binary -BinaryFile $binaryFile
    
    # Test installation
    $installedPath = Join-Path $INSTALL_DIR "$BINARY_NAME.exe"
    if (Test-Path $installedPath) {
        Write-Info "Installation successful!"
        Write-Host ""
        
        # Try to run version command
        try {
            & $installedPath --version
        }
        catch {
            Write-Info "Run the following command to verify:"
            Write-Host "  $installedPath --version"
        }
    }
    
    Write-Host ""
    Write-Host "To use with MCP, add the server to your MCP configuration."
    Write-Host "For more information, visit: https://github.com/$REPO"
}

# Run main function
Main
