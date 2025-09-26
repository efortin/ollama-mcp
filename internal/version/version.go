package version

import (
	"fmt"
	"runtime"
)

// These variables are set by GoReleaser at build time
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// String returns the version information as a formatted string
func String() string {
	return fmt.Sprintf("ollama-mcp %s (%s) built on %s for %s/%s",
		Version, Commit, Date, runtime.GOOS, runtime.GOARCH)
}

// Short returns just the version number
func Short() string {
	return Version
}
