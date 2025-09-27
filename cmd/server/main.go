package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/efortin/ollama-mcp/internal/core"
	"github.com/efortin/ollama-mcp/internal/version"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	// Parse command line flags
	versionFlag := flag.Bool("version", false, "Print version information")
	hostFlag := flag.String("host", "", "Ollama host URL (e.g., https://ollama.empyr.cloud)")
	contextSizeFlag := flag.Int("context-size", core.DefaultContextSize, "Context size for models")
	codeModelFlag := flag.String("code-model", core.DefaultCodeModel, "Model to use for code generation")
	chatModelFlag := flag.String("chat-model", core.DefaultChatModel, "Model to use for chat")
	keepAliveFlag := flag.String("keep-alive", core.DefaultKeepAlive, "Keep-alive duration for models")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println(version.String())
		os.Exit(0)
	}

	// Load configuration with command-line arguments
	config, err := core.LoadConfigFromFlags(*hostFlag, *contextSizeFlag, *codeModelFlag, *chatModelFlag, *keepAliveFlag)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create our server instance with dependency injection
	ollamaServer := core.NewServer(config)
	handlerFactory := core.NewHandlerFactory(ollamaServer)

	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{Name: "ollama-mcp", Version: version.Short()}, nil)

	// Get model names for descriptions
	codeModel := config.CodeModel
	chatModel := config.ChatModel

	// Add the code tool with the configured model
	mcp.AddTool(server, &mcp.Tool{Name: "code", Description: fmt.Sprintf("code with %s", codeModel)}, handlerFactory.CodeHandler())

	// Add the chat tool with the configured model
	mcp.AddTool(server, &mcp.Tool{Name: "chat", Description: fmt.Sprintf("chat with %s", chatModel)}, handlerFactory.ChatHandler())

	// Add the list models tool
	mcp.AddTool(server, &mcp.Tool{Name: "list-models", Description: "list available Ollama models"}, handlerFactory.ListModelsHandler())

	// Add the model info tool
	mcp.AddTool(server, &mcp.Tool{Name: "model-info", Description: "get information about a specific Ollama model"}, handlerFactory.ModelInfoHandler())

	// Add the pull model tool
	mcp.AddTool(server, &mcp.Tool{Name: "pull-model", Description: "pull a model from the Ollama library"}, handlerFactory.PullModelHandler())

	// Run the server over stdin/stdout, until the client disconnects
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
