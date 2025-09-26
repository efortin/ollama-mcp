package core

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/ollama/ollama/api"
)

// Default configuration values
const (
	DefaultContextSize = 32000
	DefaultCodeModel   = "qwen3-coder:30b"
	DefaultChatModel   = "gpt-oss:20b"
	DefaultKeepAlive   = "1m"
)

// Config holds the configuration for the Ollama MCP server
type Config struct {
	Client      *api.Client
	ContextSize int
	CodeModel   string
	ChatModel   string
	KeepAlive   string
}

// LoadConfig creates a new configuration from environment variables
func LoadConfig() (*Config, error) {
	// Use the official client's environment-based constructor
	// This respects OLLAMA_HOST environment variable
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	// If OLLAMA_HOST is not set and ClientFromEnvironment uses a default,
	// we can optionally override with a custom HTTP client
	if os.Getenv("OLLAMA_HOST") == "" && os.Getenv("OLLAMA_CUSTOM_CLIENT") == "true" {
		baseURL, _ := url.Parse("http://localhost:11434")
		client = api.NewClient(baseURL, createHTTPClient())
	}

	// Get context size from environment or use default
	contextSize := DefaultContextSize
	if sizeStr := os.Getenv("OLLAMA_CONTEXT_SIZE"); sizeStr != "" {
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 {
			contextSize = size
		}
	}

	return &Config{
		Client:      client,
		ContextSize: contextSize,
		CodeModel:   getEnvOrDefault("OLLAMA_CODE_MODEL", DefaultCodeModel),
		ChatModel:   getEnvOrDefault("OLLAMA_CHAT_MODEL", DefaultChatModel),
		KeepAlive:   getEnvOrDefault("OLLAMA_KEEP_ALIVE", DefaultKeepAlive),
	}, nil
}

// GetModel returns the model for the specified tool
func (c *Config) GetModel(toolName string) string {
	switch toolName {
	case "code":
		return c.CodeModel
	case "chat":
		return c.ChatModel
	default:
		return c.ChatModel
	}
}

// Global config instance (loaded once at startup)
var globalConfig *Config

// InitConfig initializes the global configuration
func InitConfig() error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}
	globalConfig = config
	return nil
}

// GetConfig returns the global configuration
func GetConfig() *Config {
	if globalConfig == nil {
		// Fallback: create config on demand if not initialized
		config, _ := LoadConfig()
		globalConfig = config
	}
	return globalConfig
}

// GetClient returns the Ollama client
func GetClient() *api.Client {
	return GetConfig().Client
}

// GetDefaultModel returns the default model for a tool
func GetDefaultModel(toolName string) string {
	return GetConfig().GetModel(toolName)
}

// GetDefaultContextSize returns the default context size
func GetDefaultContextSize() int {
	return GetConfig().ContextSize
}

// GetDefaultKeepAlive returns the default keep-alive duration
func GetDefaultKeepAlive() string {
	return GetConfig().KeepAlive
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// createHTTPClient creates an HTTP client with custom settings
func createHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       1 * time.Minute,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 5 * time.Minute,
	}

	return &http.Client{
		Transport: transport,
		Timeout:   10 * time.Minute,
	}
}
