package main_test

import (
	"flag"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/efortin/ollama-mcp/internal/version"
)

var _ = Describe("Main", func() {
	var (
		oldArgs []string
	)

	BeforeEach(func() {
		// Save original command line args
		oldArgs = os.Args
		// Reset flag package state
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	})

	AfterEach(func() {
		// Restore original args
		os.Args = oldArgs
	})

	Describe("Flag parsing", func() {
		Context("with version flag", func() {
			BeforeEach(func() {
				os.Args = []string{"cmd", "-version"}
			})

			It("should parse version flag correctly", func() {
				versionFlag := flag.Bool("version", false, "Print version information")
				flag.Parse()

				Expect(*versionFlag).To(BeTrue())
			})
		})

		Context("with no flags", func() {
			BeforeEach(func() {
				os.Args = []string{"cmd"}
			})

			It("should have version flag as false", func() {
				versionFlag := flag.Bool("version", false, "Print version information")
				flag.Parse()

				Expect(*versionFlag).To(BeFalse())
			})
		})
	})

	Describe("Version information", func() {
		It("should provide version string", func() {
			versionStr := version.String()
			Expect(versionStr).NotTo(BeEmpty())
		})

		It("should provide short version", func() {
			shortVersion := version.Short()
			Expect(shortVersion).NotTo(BeEmpty())
		})

		It("should have consistent version information", func() {
			versionStr := version.String()
			shortVersion := version.Short()

			Expect(versionStr).To(ContainSubstring("ollama-mcp"))
			Expect(shortVersion).NotTo(BeEmpty())
		})
	})

	Describe("Main components", func() {
		It("should be able to initialize flag parsing", func() {
			flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
			versionFlag := flag.Bool("version", false, "Print version information")

			Expect(versionFlag).NotTo(BeNil())
		})

		It("should have access to version information", func() {
			// Verify that the main function components can access version info
			versionStr := version.String()
			shortVersion := version.Short()

			Expect(versionStr).NotTo(BeEmpty())
			Expect(shortVersion).NotTo(BeEmpty())
		})
	})
})
