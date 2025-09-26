package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"mcp-hello/internal/core"
	"mcp-hello/internal/version"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println(version.String())
		os.Exit(0)
	}

	// Initialize configuration
	if err := core.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Create a server with tools
	server := mcp.NewServer(&mcp.Implementation{Name: "ollama-mcp", Version: version.Short()}, nil)

	// Get model names from environment or use defaults
	codeModel := os.Getenv("OLLAMA_CODE_MODEL")
	if codeModel == "" {
		codeModel = "qwen3-coder:30b"
	}

	chatModel := os.Getenv("OLLAMA_CHAT_MODEL")
	if chatModel == "" {
		chatModel = "gpt-oss:20b"
	}

	// Add the code tool with the configured model
	mcp.AddTool(server, &mcp.Tool{Name: "code", Description: fmt.Sprintf("code with %s", codeModel)}, core.ChatWithOllama)

	// Add the chat tool with the configured model
	mcp.AddTool(server, &mcp.Tool{Name: "chat", Description: fmt.Sprintf("chat with %s", chatModel)}, core.ChatWithOllama)

	// Add the list models tool
	mcp.AddTool(server, &mcp.Tool{Name: "list-models", Description: "list available Ollama models"}, core.ListModels)

	// Add the model info tool
	mcp.AddTool(server, &mcp.Tool{Name: "model-info", Description: "get information about a specific Ollama model"}, core.ModelInfo)

	// Add the pull model tool
	mcp.AddTool(server, &mcp.Tool{Name: "pull-model", Description: "pull a model from the Ollama library"}, core.PullModel)

	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
