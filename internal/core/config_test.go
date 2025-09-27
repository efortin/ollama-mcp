package core_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/efortin/ollama-mcp/internal/core"
)

var _ = Describe("Config", func() {
	var (
		originalHost        string
		originalContextSize string
		originalCodeModel   string
		originalChatModel   string
		originalKeepAlive   string
	)

	BeforeEach(func() {
		// Save original environment
		originalHost = os.Getenv("OLLAMA_HOST")
		originalContextSize = os.Getenv("OLLAMA_CONTEXT_SIZE")
		originalCodeModel = os.Getenv("OLLAMA_CODE_MODEL")
		originalChatModel = os.Getenv("OLLAMA_CHAT_MODEL")
		originalKeepAlive = os.Getenv("OLLAMA_KEEP_ALIVE")
	})

	AfterEach(func() {
		// Restore original environment
		_ = os.Setenv("OLLAMA_HOST", originalHost)
		_ = os.Setenv("OLLAMA_CONTEXT_SIZE", originalContextSize)
		_ = os.Setenv("OLLAMA_CODE_MODEL", originalCodeModel)
		_ = os.Setenv("OLLAMA_CHAT_MODEL", originalChatModel)
		_ = os.Setenv("OLLAMA_KEEP_ALIVE", originalKeepAlive)
	})

	Describe("LoadConfig", func() {
		Context("with default values", func() {
			BeforeEach(func() {
				// Clear environment variables to test defaults
				_ = os.Unsetenv("OLLAMA_HOST")
				_ = os.Unsetenv("OLLAMA_CONTEXT_SIZE")
				_ = os.Unsetenv("OLLAMA_CODE_MODEL")
				_ = os.Unsetenv("OLLAMA_CHAT_MODEL")
				_ = os.Unsetenv("OLLAMA_KEEP_ALIVE")
			})

			It("should load configuration with default values", func() {
				config, err := core.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config.ContextSize).To(Equal(core.DefaultContextSize))
				Expect(config.CodeModel).To(Equal(core.DefaultCodeModel))
				Expect(config.ChatModel).To(Equal(core.DefaultChatModel))
				Expect(config.KeepAlive).To(Equal(core.DefaultKeepAlive))
				Expect(config.Client).NotTo(BeNil())
			})
		})

		Context("with environment variables", func() {
			BeforeEach(func() {
				_ = os.Setenv("OLLAMA_CONTEXT_SIZE", "16000")
				_ = os.Setenv("OLLAMA_CODE_MODEL", "test-code-model")
				_ = os.Setenv("OLLAMA_CHAT_MODEL", "test-chat-model")
				_ = os.Setenv("OLLAMA_KEEP_ALIVE", "5m")
			})

			It("should load configuration from environment variables", func() {
				config, err := core.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config.ContextSize).To(Equal(16000))
				Expect(config.CodeModel).To(Equal("test-code-model"))
				Expect(config.ChatModel).To(Equal("test-chat-model"))
				Expect(config.KeepAlive).To(Equal("5m"))
			})
		})

		Context("with invalid context size", func() {
			BeforeEach(func() {
				_ = os.Setenv("OLLAMA_CONTEXT_SIZE", "invalid")
			})

			It("should fallback to default when invalid", func() {
				config, err := core.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config.ContextSize).To(Equal(core.DefaultContextSize))
			})
		})

		Context("with zero context size", func() {
			BeforeEach(func() {
				_ = os.Setenv("OLLAMA_CONTEXT_SIZE", "0")
			})

			It("should fallback to default when zero", func() {
				config, err := core.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config.ContextSize).To(Equal(core.DefaultContextSize))
			})
		})
	})

	Describe("Config.GetModel", func() {
		var config *core.Config

		BeforeEach(func() {
			config = &core.Config{
				CodeModel: "test-code-model",
				ChatModel: "test-chat-model",
			}
		})

		DescribeTable("tool name to model mapping",
			func(toolName, expectedModel string) {
				model := config.GetModel(toolName)
				Expect(model).To(Equal(expectedModel))
			},
			Entry("code tool", "code", "test-code-model"),
			Entry("chat tool", "chat", "test-chat-model"),
			Entry("unknown tool", "unknown", "test-chat-model"),
			Entry("empty tool name", "", "test-chat-model"),
		)
	})

	Describe("Server", func() {
		var config *core.Config

		BeforeEach(func() {
			config = &core.Config{
				ContextSize: 16000,
				CodeModel:   "test-code",
				ChatModel:   "test-chat",
				KeepAlive:   "2m",
			}
		})

		Describe("NewServer", func() {
			It("should create a new server", func() {
				server := core.NewServer(config)
				Expect(server).NotTo(BeNil())
				Expect(server.GetConfig()).To(Equal(config))
			})
		})

		Describe("Server methods", func() {
			var server *core.Server

			BeforeEach(func() {
				server = core.NewServer(config)
			})

			It("should return the correct config", func() {
				returnedConfig := server.GetConfig()
				Expect(returnedConfig).To(Equal(config))
			})

			It("should return the correct client", func() {
				client := server.GetClient()
				Expect(client).To(Equal(config.Client))
			})

			DescribeTable("default model retrieval",
				func(toolName, expectedModel string) {
					model := server.GetDefaultModel(toolName)
					Expect(model).To(Equal(expectedModel))
				},
				Entry("code tool", "code", "test-code"),
				Entry("chat tool", "chat", "test-chat"),
				Entry("unknown tool", "unknown", "test-chat"),
			)

			It("should return the correct default context size", func() {
				contextSize := server.GetDefaultContextSize()
				Expect(contextSize).To(Equal(16000))
			})

			It("should return the correct default keep alive", func() {
				keepAlive := server.GetDefaultKeepAlive()
				Expect(keepAlive).To(Equal("2m"))
			})
		})
	})

	Describe("Helper functions", func() {
		Describe("getEnvOrDefault", func() {
			const testKey = "TEST_ENV_VAR_UNIQUE"

			AfterEach(func() {
				_ = os.Unsetenv(testKey)
			})

			Context("with unset environment variable", func() {
				BeforeEach(func() {
					_ = os.Unsetenv(testKey)
				})

				It("should return default value", func() {
					// We can't test the private function directly, but we can test through LoadConfig
					// This is more of an integration test
					_ = os.Unsetenv("OLLAMA_CODE_MODEL")
					config, err := core.LoadConfig()
					Expect(err).NotTo(HaveOccurred())
					Expect(config.CodeModel).To(Equal(core.DefaultCodeModel))
				})
			})

			Context("with set environment variable", func() {
				BeforeEach(func() {
					_ = os.Setenv("OLLAMA_CODE_MODEL", "custom-model")
				})

				It("should return environment value", func() {
					config, err := core.LoadConfig()
					Expect(err).NotTo(HaveOccurred())
					Expect(config.CodeModel).To(Equal("custom-model"))
				})
			})
		})

		Describe("createHTTPClient", func() {
			It("should be tested through LoadConfig", func() {
				config, err := core.LoadConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(config.Client).NotTo(BeNil())
			})
		})
	})
})
