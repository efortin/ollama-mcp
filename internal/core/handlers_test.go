package core_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/efortin/ollama-mcp/internal/core"
)

var _ = Describe("Handlers", func() {
	var (
		server  *core.Server
		factory *core.HandlerFactory
	)

	BeforeEach(func() {
		config := &core.Config{
			Client:      nil, // We'll test validation logic that doesn't require actual client
			ContextSize: 32000,
			CodeModel:   "test-code-model",
			ChatModel:   "test-chat-model",
			KeepAlive:   "1m",
		}
		server = core.NewServer(config)
		factory = core.NewHandlerFactory(server)
	})

	Describe("NewHandlerFactory", func() {
		It("should create a new handler factory", func() {
			Expect(factory).NotTo(BeNil())
			Expect(factory.GetServer()).To(Equal(server))
		})
	})

	Describe("Handler Creation", func() {
		It("should create all handlers without panicking", func() {
			chatHandler := factory.ChatHandler()
			Expect(chatHandler).NotTo(BeNil())

			codeHandler := factory.CodeHandler()
			Expect(codeHandler).NotTo(BeNil())

			listHandler := factory.ListModelsHandler()
			Expect(listHandler).NotTo(BeNil())

			infoHandler := factory.ModelInfoHandler()
			Expect(infoHandler).NotTo(BeNil())

			pullHandler := factory.PullModelHandler()
			Expect(pullHandler).NotTo(BeNil())
		})
	})

	Describe("ValidateModelName", func() {
		DescribeTable("model name validation",
			func(modelName string, shouldBeValid bool) {
				err := factory.ValidateModelName(modelName)
				if shouldBeValid {
					Expect(err).To(BeNil())
				} else {
					Expect(err).To(HaveOccurred())
				}
			},
			Entry("valid model name", "llama2:7b", true),
			Entry("empty model name", "", false),
			Entry("path traversal", "../malicious", false),
			Entry("forward slash", "model/name", false),
			Entry("suspicious chars", "model<script>", false),
		)
	})

	Describe("ValidateChatInput", func() {
		DescribeTable("chat input validation",
			func(input core.ChatInput, shouldBeValid bool) {
				err := factory.ValidateChatInput(input)
				if shouldBeValid {
					Expect(err).To(BeNil())
				} else {
					Expect(err).To(HaveOccurred())
				}
			},
			Entry("valid input", core.ChatInput{Message: "Hello"}, true),
			Entry("empty message", core.ChatInput{Message: ""}, false),
			Entry("negative temperature", core.ChatInput{
				Message:     "Hello",
				Temperature: func() *float32 { v := float32(-0.1); return &v }(),
			}, false),
			Entry("high temperature", core.ChatInput{
				Message:     "Hello",
				Temperature: func() *float32 { v := float32(2.1); return &v }(),
			}, false),
			Entry("negative top_p", core.ChatInput{
				Message: "Hello",
				TopP:    func() *float32 { v := float32(-0.1); return &v }(),
			}, false),
			Entry("high top_p", core.ChatInput{
				Message: "Hello",
				TopP:    func() *float32 { v := float32(1.1); return &v }(),
			}, false),
			Entry("negative top_k", core.ChatInput{
				Message: "Hello",
				TopK:    func() *int { v := -1; return &v }(),
			}, false),
		)
	})
})
