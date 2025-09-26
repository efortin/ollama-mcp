# Ollama MCP

A Go project that implements a Model Context Protocol (MCP) server with Ollama integration for chat functionality and model management.

## Quick Install

### macOS/Linux

```bash
curl -fsSL https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.sh | bash
```

### Windows (PowerShell)

```powershell
iwr -useb https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.ps1 | iex
```

### Manual Download

You can also download the binary directly from the [releases page](https://github.com/efortin/ollama-mcp/releases/latest).

### Installation Options

The install scripts support the following environment variables:

- `OLLAMA_MCP_INSTALL_DIR`: Custom installation directory (default: `~/.local/bin` on Unix, `%LOCALAPPDATA%\Programs\ollama-mcp` on Windows)
- `OLLAMA_MCP_VERSION`: Install a specific version (default: latest)

Example:
```bash
# Install to a custom directory
OLLAMA_MCP_INSTALL_DIR=/opt/ollama-mcp curl -fsSL https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.sh | bash

# Install a specific version
OLLAMA_MCP_VERSION=v1.0.0 curl -fsSL https://raw.githubusercontent.com/efortin/ollama-mcp/main/docs/install.sh | bash
```

## Project Description

This project provides an MCP server that exposes several tools for interacting with Ollama:

1. **Code Tool**: Code generation and assistance (uses model from OLLAMA_CODE_MODEL or defaults to qwen3-coder:30b)
2. **Chat Tool**: General conversations (uses model from OLLAMA_CHAT_MODEL or defaults to gpt-oss:20b)
3. **List Models Tool**: List all available Ollama models
4. **Model Info Tool**: Get detailed information about a specific Ollama model
5. **Pull Model Tool**: Pull a model from the Ollama library

The project follows the Go standard project layout and uses mise for environment management.

## Prerequisites

- [mise](https://mise.jdx.dev/) installed
- Go 1.21+ (managed by mise)
- [Ollama](https://ollama.ai/) running locally on port 11434 (or set OLLAMA_HOST environment variable)
- Git

## Quickstart

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/ollama-mcp.git
cd ollama-mcp

# Install dependencies
mise run setup

# Build the server
mise run build

# Install binary (optional)
cp bin/server /usr/local/bin/ollama-mcp
```

### Quick Start

```bash
# Run the server directly
mise run run

# Or run the installed binary
ollama-mcp
```

## Project Structure

```
ollama-mcp/
├── cmd/
│   └── server/         # Application entrypoint
│       └── main.go
├── internal/
│   └── core/           # Internal application code
│       ├── chat.go     # Chat functionality using Ollama
│       ├── code.go     # Code generation functionality
│       ├── config.go   # Configuration for Ollama client
│       ├── models.go   # Model listing and info functionality
│       └── pull.go     # Model pulling functionality
├── pkg/                # Library code that can be used by external applications
├── tests/              # Test files
│   ├── chat_test.go    # Tests for chat functionality
│   └── ollama_test.go  # Tests for Ollama functionality
├── .mise.toml          # mise configuration
├── .mise-tasks.toml    # mise tasks
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
└── README.md           # This file
```

## Development

### Building and Installing the Binary

To compile and install the MCP server binary, always use mise commands:

```bash
# Build the server
mise run build

# The binary will be created in bin/server
# To install it to your system, you can copy it manually:
cp bin/server /usr/local/bin/ollama-mcp

# Or use mise to build for all platforms
mise run build-all
```

For development builds:

```bash
# Run the server directly (builds and runs)
mise run run

# Run tests
mise run test

# Run tests with verbose output
mise run test-verbose

# Run linter
mise run lint
```

### Available Tasks

- `mise run setup` - Set up dependencies
- `mise run build` - Build the server
- `mise run run` - Run the server
- `mise run test` - Run tests
- `mise run test-verbose` - Run tests with verbose output
- `mise run lint` - Run linter
- `mise run build-all` - Build for all platforms
- `mise run deps` - Update dependencies
- `mise run update-deps` - Update dependencies to latest versions

### Cross-compilation

The project supports cross-compilation for the following platforms:
- macOS (amd64, arm64)
- Linux (amd64, arm64)

Use `mise run build-all` to build for all platforms.

## Configuration

The Ollama client can be configured using the following environment variables:

- `OLLAMA_HOST`: The URL of the Ollama server (default: http://localhost:11434)
- `OLLAMA_CONTEXT_SIZE`: Maximum context size in tokens (default: 32000)
- `OLLAMA_CODE_MODEL`: Model for code tool (default: qwen3-coder:30b)
- `OLLAMA_CHAT_MODEL`: Model for chat tool (default: gpt-oss:20b)
- `OLLAMA_KEEP_ALIVE`: Duration to keep models loaded in VRAM (default: 1m)

## Troubleshooting

### Error: "invalid character '<' looking for beginning of value"

This error indicates that the Ollama server is returning HTML instead of JSON. Common causes and solutions:

1. **Ollama is not running**
   ```bash
   # Start Ollama
   ollama serve
   ```

2. **Wrong port or host**
   ```bash
   # Check if Ollama is running on the expected port
   curl http://localhost:11434/api/tags
   
   # If using a different host/port, set the environment variable
   export OLLAMA_HOST=http://your-host:your-port
   ```

3. **Check Ollama connection**
   ```bash
   # Test if Ollama is accessible
   curl http://localhost:11434/api/tags
   ```

### Common Issues

- **Models not found**: Make sure you have pulled the required models:
  ```bash
  ollama pull qwen3-coder:30b
  ollama pull gpt-oss:20b
  # Or pull your preferred models and set environment variables
  export OLLAMA_CODE_MODEL=your-code-model
  export OLLAMA_CHAT_MODEL=your-chat-model
  ```

- **Connection refused**: Ensure Ollama is running and accessible
- **Timeout errors**: Check if Ollama is responding slowly due to model loading

## Available Tools

### Code Tool

Code with qwen3-coder:30b model for code generation, debugging, and programming assistance. The model stays loaded in VRAM based on the keep-alive duration (default: 1 minute) to improve performance for consecutive requests.

### Chat Tool

Chat with gpt-oss:20b model for general conversations and text generation. The model stays loaded in VRAM based on the keep-alive duration (default: 1 minute) to improve performance for consecutive requests.

### List Models Tool

List all available Ollama models.

### Model Info Tool

Get detailed information about a specific Ollama model.

### Pull Model Tool

Pull a model from the Ollama library.

## Using the MCP Server

Once the server is running, you can interact with it through any MCP-compatible client. The server exposes the following tools:

- **chat**: General conversations with AI models
- **code**: Code generation and programming assistance
- **list-models**: List all available Ollama models
- **model-info**: Get detailed information about a specific model
- **pull-model**: Download models from the Ollama library

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

This means:
- ✅ You can use, modify, and distribute this software
- ✅ You must provide source code when distributing
- ⚠️ If you use this in a network service, you must provide source code to users
- ⚠️ Any modifications must also be licensed under AGPL-3.0

For commercial use with different licensing terms, please contact me.

See the [LICENSE](LICENSE) file for full details.
