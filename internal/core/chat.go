package core

// ChatInput represents the input for chat operations
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

// Note: ChatWithOllama is deprecated. Use HandlerFactory.ChatHandler() instead.
// This function is kept for backward compatibility but should not be used directly.
