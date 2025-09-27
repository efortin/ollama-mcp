package version_test

import (
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/efortin/ollama-mcp/internal/version"
)

var _ = Describe("Version", func() {

	Describe("String", func() {
		Context("with test values", func() {
			It("should return formatted version string", func() {
				result := version.String()

				expectedParts := []string{
					"ollama-mcp",
					runtime.GOOS,
					runtime.GOARCH,
				}

				for _, part := range expectedParts {
					Expect(result).To(ContainSubstring(part))
				}
			})

			It("should not be empty", func() {
				result := version.String()
				Expect(result).NotTo(BeEmpty())
			})
		})
	})

	Describe("Short", func() {
		It("should return version number", func() {
			result := version.Short()
			Expect(result).NotTo(BeEmpty())
		})

		It("should be consistent", func() {
			result1 := version.Short()
			result2 := version.Short()
			Expect(result1).To(Equal(result2))
		})
	})

	Describe("Default values", func() {
		It("should have default values set", func() {
			// Test that default values are not empty
			versionStr := version.String()
			Expect(versionStr).NotTo(BeEmpty())

			shortVersion := version.Short()
			Expect(shortVersion).NotTo(BeEmpty())
		})

		It("should contain expected format", func() {
			versionStr := version.String()
			// Should follow format: "ollama-mcp VERSION (COMMIT) built on DATE for GOOS/GOARCH"
			Expect(versionStr).To(ContainSubstring("ollama-mcp"))
			Expect(versionStr).To(ContainSubstring("built on"))
			Expect(versionStr).To(ContainSubstring("for"))
		})
	})
})
