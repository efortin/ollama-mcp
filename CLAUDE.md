# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

This project uses [mise](https://mise.jdx.dev/) for task management. All development commands should be run through mise:

### Core Development Tasks
- `mise run build` - Build the server binary to bin/server
- `mise run run` - Run the server directly (builds and runs)
- `mise run test` - Run all tests
- `mise run test-verbose` - Run tests with verbose output
- `mise run lint` - Run golangci-lint

### Setup and Dependencies
- `mise run setup` - Initial setup (run go mod tidy and update deps)
- `mise run deps` - Update go.mod dependencies
- `mise run update-deps` - Update all dependencies to latest versions

### Cross-compilation
- `mise run build-all` - Build for all supported platforms
- Individual platform builds: `build-darwin-amd64`, `build-darwin-arm64`, `build-linux-amd64`, `build-linux-arm64`

### Testing Individual Components
The project uses **Ginkgo BDD testing framework**. To test specific packages:
```bash
# Run tests for specific packages
go test ./internal/core -v          # Core functionality tests
go test ./internal/version -v       # Version tests
go test ./cmd/server -v             # Main server tests

# Run tests with coverage
go test ./... -cover

# Generate detailed coverage report
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

# Run legacy integration tests
go test ./tests -v
```

## Architecture Overview

This is a **Model Context Protocol (MCP) server** that provides Ollama integration for AI chat and code generation tools.

### Core Architecture Pattern
The codebase uses **dependency injection** with a clean separation of concerns:

1. **Configuration Layer** (`internal/core/config.go`): Handles environment-based configuration and Ollama client creation
2. **Server Layer** (`internal/core/config.go:Server`): Wrapper that holds configuration and provides access methods
3. **Handler Factory** (`internal/core/handlers.go`): Creates tool handlers with injected dependencies
4. **Tool Handlers**: Individual functions for each MCP tool (chat, code, list-models, etc.)

### Key Components

**Main Entry Point** (`cmd/server/main.go`):
- Parses command line flags (--version)
- Loads configuration from environment
- Creates server instance with dependency injection
- Registers MCP tools using the handler factory
- Runs MCP server over stdin/stdout

**Configuration System** (`internal/core/config.go`):
- Uses official Ollama API client with environment-based configuration
- Supports OLLAMA_HOST, OLLAMA_CONTEXT_SIZE, OLLAMA_CODE_MODEL, OLLAMA_CHAT_MODEL, OLLAMA_KEEP_ALIVE
- Provides default values: qwen3-coder:30b for code, gpt-oss:20b for chat
- Creates HTTP client with custom timeouts and connection pooling

**Tool Implementations**:
- **Chat Tool**: General conversations using configured chat model
- **Code Tool**: Code assistance using configured code model (with default system prompt)
- **List Models**: Lists available Ollama models
- **Model Info**: Gets detailed information about specific models
- **Pull Model**: Downloads models from Ollama library

### Security Features
- Input validation for all user inputs
- Model name validation prevents path traversal attacks
- Timeout contexts for all Ollama API calls
- Clean separation between configuration and business logic

### MCP Integration
Uses the official `github.com/modelcontextprotocol/go-sdk/mcp` package for:
- Tool registration and handling
- Stdin/stdout transport for client communication
- Request/response serialization

## Environment Configuration

The server respects standard Ollama environment variables:
- `OLLAMA_HOST` - Ollama server URL (default: http://localhost:11434)
- `OLLAMA_CONTEXT_SIZE` - Token context size (default: 32000)
- `OLLAMA_CODE_MODEL` - Model for code tool (default: qwen3-coder:30b)
- `OLLAMA_CHAT_MODEL` - Model for chat tool (default: gpt-oss:20b)
- `OLLAMA_KEEP_ALIVE` - Model keep-alive duration (default: 1m)

## Testing Strategy

The project uses **Ginkgo BDD testing framework** with **Gomega matchers** for readable, well-organized tests:

### Test Structure
- **Core functionality**: `internal/core/*_test.go` - BDD tests for handlers, configuration, and validation
- **Version module**: `internal/version/*_test.go` - Version string formatting tests
- **Main server**: `cmd/server/*_test.go` - Integration tests for CLI and main components
- **Legacy tests**: `tests/*.go` - Legacy integration tests (maintained for compatibility)

### Test Organization Pattern
Tests follow Ginkgo's hierarchical structure:
```go
Describe("ComponentName", func() {
    Context("when condition", func() {
        It("should behavior", func() {
            // Test implementation with Gomega matchers
        })
    })
})
```

### Coverage Goals
- **Current coverage: 30.4%** overall
- `internal/version`: **100%** coverage
- `internal/core`: **33.8%** coverage
- Focus on validation functions and configuration logic

Always run `mise run test` to execute the full test suite before making changes.

## Project Structure

The codebase follows Go standard project layout with these key differences from the README:

```
ollama-mcp/
├── cmd/server/                 # Application entrypoint
│   ├── main.go                # CLI parsing and MCP server setup
│   ├── main_test.go           # Ginkgo tests for main functionality
│   └── main_suite_test.go     # Ginkgo test suite
├── internal/core/             # Core business logic (private)
│   ├── config.go              # Environment configuration and HTTP client
│   ├── handlers.go            # MCP tool handlers with dependency injection
│   ├── chat.go                # Chat input/output types
│   ├── code.go                # Code input/output types
│   ├── models.go              # Model listing types
│   ├── pull.go                # Model pulling types
│   ├── *_test.go              # Ginkgo BDD tests
│   └── core_suite_test.go     # Ginkgo test suite
├── internal/version/          # Version information module
│   ├── version.go             # Version string formatting
│   ├── version_test.go        # Ginkgo tests
│   └── version_suite_test.go  # Ginkgo test suite
└── tests/                     # Legacy integration tests
    └── *.go                   # Original test files (maintained for compatibility)
```

**Note**: The README's project structure is outdated. The actual implementation separates types into individual files (chat.go, code.go, models.go, pull.go) rather than having combined functionality files.