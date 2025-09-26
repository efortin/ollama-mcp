package tests

import (
	"os"
	"testing"

	"github.com/efortin/ollama-mcp/internal/core"
)

func TestLoadConfig(t *testing.T) {
	// Save original environment
	originalHost := os.Getenv("OLLAMA_HOST")
	originalContextSize := os.Getenv("OLLAMA_CONTEXT_SIZE")
	originalCodeModel := os.Getenv("OLLAMA_CODE_MODEL")
	originalChatModel := os.Getenv("OLLAMA_CHAT_MODEL")
	originalKeepAlive := os.Getenv("OLLAMA_KEEP_ALIVE")
	defer func() {
		_ = os.Setenv("OLLAMA_HOST", originalHost)
		_ = os.Setenv("OLLAMA_CONTEXT_SIZE", originalContextSize)
		_ = os.Setenv("OLLAMA_CODE_MODEL", originalCodeModel)
		_ = os.Setenv("OLLAMA_CHAT_MODEL", originalChatModel)
		_ = os.Setenv("OLLAMA_KEEP_ALIVE", originalKeepAlive)
	}()

	// Clear environment variables to test defaults
	_ = os.Unsetenv("OLLAMA_HOST")
	_ = os.Unsetenv("OLLAMA_CONTEXT_SIZE")
	_ = os.Unsetenv("OLLAMA_CODE_MODEL")
	_ = os.Unsetenv("OLLAMA_CHAT_MODEL")
	_ = os.Unsetenv("OLLAMA_KEEP_ALIVE")

	// Test with default values
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if config.ContextSize != core.DefaultContextSize {
		t.Errorf("LoadConfig() ContextSize = %v, want %v", config.ContextSize, core.DefaultContextSize)
	}

	if config.CodeModel != core.DefaultCodeModel {
		t.Errorf("LoadConfig() CodeModel = %v, want %v", config.CodeModel, core.DefaultCodeModel)
	}

	if config.ChatModel != core.DefaultChatModel {
		t.Errorf("LoadConfig() ChatModel = %v, want %v", config.ChatModel, core.DefaultChatModel)
	}

	if config.KeepAlive != core.DefaultKeepAlive {
		t.Errorf("LoadConfig() KeepAlive = %v, want %v", config.KeepAlive, core.DefaultKeepAlive)
	}
}

func TestLoadConfig_WithEnvVars(t *testing.T) {
	// Save original environment
	originalHost := os.Getenv("OLLAMA_HOST")
	originalContextSize := os.Getenv("OLLAMA_CONTEXT_SIZE")
	originalCodeModel := os.Getenv("OLLAMA_CODE_MODEL")
	originalChatModel := os.Getenv("OLLAMA_CHAT_MODEL")
	originalKeepAlive := os.Getenv("OLLAMA_KEEP_ALIVE")
	defer func() {
		_ = os.Setenv("OLLAMA_HOST", originalHost)
		_ = os.Setenv("OLLAMA_CONTEXT_SIZE", originalContextSize)
		_ = os.Setenv("OLLAMA_CODE_MODEL", originalCodeModel)
		_ = os.Setenv("OLLAMA_CHAT_MODEL", originalChatModel)
		_ = os.Setenv("OLLAMA_KEEP_ALIVE", originalKeepAlive)
	}()

	// Set test environment variables
	testHost := "http://test-host:8080"
	testContextSize := "16000"
	testCodeModel := "test-code-model"
	testChatModel := "test-chat-model"
	testKeepAlive := "5m"

	_ = os.Setenv("OLLAMA_HOST", testHost)
	_ = os.Setenv("OLLAMA_CONTEXT_SIZE", testContextSize)
	_ = os.Setenv("OLLAMA_CODE_MODEL", testCodeModel)
	_ = os.Setenv("OLLAMA_CHAT_MODEL", testChatModel)
	_ = os.Setenv("OLLAMA_KEEP_ALIVE", testKeepAlive)

	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if config.ContextSize != 16000 {
		t.Errorf("LoadConfig() ContextSize = %v, want %v", config.ContextSize, 16000)
	}

	if config.CodeModel != testCodeModel {
		t.Errorf("LoadConfig() CodeModel = %v, want %v", config.CodeModel, testCodeModel)
	}

	if config.ChatModel != testChatModel {
		t.Errorf("LoadConfig() ChatModel = %v, want %v", config.ChatModel, testChatModel)
	}

	if config.KeepAlive != testKeepAlive {
		t.Errorf("LoadConfig() KeepAlive = %v, want %v", config.KeepAlive, testKeepAlive)
	}
}

func TestServerConfig(t *testing.T) {
	// Load config
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Create server
	server := core.NewServer(config)
	if server == nil {
		t.Error("NewServer() returned nil")
	}

	// Get config should return non-nil
	serverConfig := server.GetConfig()
	if serverConfig == nil {
		t.Error("GetConfig() returned nil")
	}
}

func TestServerGetDefaultModel(t *testing.T) {
	// Load config
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Create server
	server := core.NewServer(config)

	tests := []struct {
		name     string
		toolName string
		want     string
	}{
		{"code tool", "code", core.DefaultCodeModel},
		{"chat tool", "chat", core.DefaultChatModel},
		{"unknown tool", "unknown", core.DefaultChatModel},
		{"empty tool", "", core.DefaultChatModel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := server.GetDefaultModel(tt.toolName)
			if model != tt.want {
				t.Errorf("GetDefaultModel(%v) = %v, want %v", tt.toolName, model, tt.want)
			}
		})
	}
}

func TestServerGetDefaultContextSize(t *testing.T) {
	// Load config
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Create server
	server := core.NewServer(config)

	contextSize := server.GetDefaultContextSize()
	if contextSize != core.DefaultContextSize {
		t.Errorf("GetDefaultContextSize() = %v, want %v", contextSize, core.DefaultContextSize)
	}
}

func TestServerGetDefaultKeepAlive(t *testing.T) {
	// Load config
	config, err := core.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	// Create server
	server := core.NewServer(config)

	keepAlive := server.GetDefaultKeepAlive()
	if keepAlive != core.DefaultKeepAlive {
		t.Errorf("GetDefaultKeepAlive() = %v, want %v", keepAlive, core.DefaultKeepAlive)
	}
}
