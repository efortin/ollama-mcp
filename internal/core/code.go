package core

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type CodeInput struct {
	Message      string         `json:"message" jsonschema:"the message to send to the model"`
	ContextSize  *int           `json:"context_size,omitempty" jsonschema:"maximum context size in tokens (optional)"`
	SystemPrompt string         `json:"system_prompt,omitempty" jsonschema:"system prompt to use (optional)"`
	Options      map[string]any `json:"options,omitempty" jsonschema:"additional model options (optional)"`
	KeepAlive    *string        `json:"keep_alive,omitempty" jsonschema:"duration to keep the model loaded in memory (optional)"`
}

type CodeOutput struct {
	Response string `json:"response" jsonschema:"the response from the model"`
}

// Code sends a message to an Ollama model optimized for code generation
func Code(ctx context.Context, req *mcp.CallToolRequest, input CodeInput) (
	*mcp.CallToolResult,
	CodeOutput,
	error,
) {
	// Convert CodeInput to ChatInput and reuse ChatWithOllama
	chatInput := ChatInput{
		Model:        "", // Always use the environment variable for code tool
		Message:      input.Message,
		ContextSize:  input.ContextSize,
		SystemPrompt: input.SystemPrompt,
		Options:      input.Options,
		KeepAlive:    input.KeepAlive,
		ToolName:     "code", // Specify that this is the code tool
	}

	// If no system prompt is provided, use a default one for code
	if chatInput.SystemPrompt == "" {
		chatInput.SystemPrompt = "You are a helpful coding assistant. Provide clear, concise, and well-commented code solutions."
	}

	// Call the unified chat function
	result, chatOutput, err := ChatWithOllama(ctx, req, chatInput)
	if err != nil {
		return nil, CodeOutput{}, err
	}

	// Convert ChatOutput to CodeOutput
	return result, CodeOutput(chatOutput), nil
}
