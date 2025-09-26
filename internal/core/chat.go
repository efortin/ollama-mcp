package core

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/ollama/ollama/api"
)

type ChatInput struct {
	Model        string         `json:"model" jsonschema:"the Ollama model to use for chat"`
	Message      string         `json:"message" jsonschema:"the message to send to the model"`
	ContextSize  *int           `json:"context_size,omitempty" jsonschema:"maximum context size in tokens (optional)"`
	Temperature  *float32       `json:"temperature,omitempty" jsonschema:"controls randomness (0.0 to 1.0, optional)"`
	TopP         *float32       `json:"top_p,omitempty" jsonschema:"controls diversity via nucleus sampling (0.0 to 1.0, optional)"`
	TopK         *int           `json:"top_k,omitempty" jsonschema:"controls diversity via top-k sampling (optional)"`
	SystemPrompt string         `json:"system_prompt,omitempty" jsonschema:"system prompt to use (optional)"`
	Options      map[string]any `json:"options,omitempty" jsonschema:"additional model options (optional)"`
	ToolName     string         `json:"tool_name,omitempty" jsonschema:"name of the tool being used (optional)"`
	KeepAlive    *string        `json:"keep_alive,omitempty" jsonschema:"duration to keep the model loaded in memory (optional)"`
}

type ChatOutput struct {
	Response string `json:"response" jsonschema:"the response from the model"`
}

// ChatWithOllama sends a message to an Ollama model and returns the response
func ChatWithOllama(ctx context.Context, req *mcp.CallToolRequest, input ChatInput) (
	*mcp.CallToolResult,
	ChatOutput,
	error,
) {
	// Get configuration
	config := GetConfig()
	if config == nil {
		return nil, ChatOutput{}, fmt.Errorf("failed to get configuration")
	}

	// Determine which tool is being used
	toolName := "chat" // Default
	if input.ToolName != "" {
		toolName = strings.ToLower(input.ToolName)
		if toolName != "code" && toolName != "chat" {
			toolName = "chat"
		}
	}

	// Use default model if not specified
	modelToUse := input.Model
	if modelToUse == "" {
		modelToUse = config.GetModel(toolName)
	}

	// Build the chat request
	chatRequest := &api.ChatRequest{
		Model: modelToUse,
		Messages: []api.Message{
			{
				Role:    "user",
				Content: input.Message,
			},
		},
		Stream:  new(bool), // Set to false (non-streaming)
		Options: make(map[string]interface{}),
	}

	// Add system prompt if provided
	if input.SystemPrompt != "" {
		systemMsg := api.Message{
			Role:    "system",
			Content: input.SystemPrompt,
		}
		chatRequest.Messages = append([]api.Message{systemMsg}, chatRequest.Messages...)
	}

	// Set context size
	if input.ContextSize != nil {
		chatRequest.Options["num_ctx"] = *input.ContextSize
	} else {
		chatRequest.Options["num_ctx"] = config.ContextSize
	}

	// Set temperature
	if input.Temperature != nil {
		chatRequest.Options["temperature"] = *input.Temperature
	}

	// Set top_p
	if input.TopP != nil {
		chatRequest.Options["top_p"] = *input.TopP
	}

	// Set top_k
	if input.TopK != nil {
		chatRequest.Options["top_k"] = *input.TopK
	}

	// Add any additional options
	if input.Options != nil {
		for k, v := range input.Options {
			chatRequest.Options[k] = v
		}
	}

	// Set keep_alive
	if input.KeepAlive != nil {
		chatRequest.Options["keep_alive"] = *input.KeepAlive
	} else {
		chatRequest.Options["keep_alive"] = config.KeepAlive
	}

	// Variable to store the final response
	var finalResponse string

	// Use the official client's Chat method
	err := config.Client.Chat(ctx, chatRequest, func(response api.ChatResponse) error {
		if response.Message.Content != "" {
			finalResponse = response.Message.Content
		}
		return nil
	})

	if err != nil {
		return nil, ChatOutput{}, fmt.Errorf("failed to chat with Ollama: %w", err)
	}

	return nil, ChatOutput{Response: finalResponse}, nil
}
