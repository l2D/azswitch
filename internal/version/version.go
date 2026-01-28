// Package version provides build version information.
package version

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time via ldflags.
var (
	Version   = "dev"
	CommitSHA = "unknown"
	BuildTime = "unknown"
)

// Info returns formatted version information.
func Info() string {
	return fmt.Sprintf("azswitch %s (%s) built %s with %s",
		Version, CommitSHA, BuildTime, runtime.Version())
}

// Short returns a short version string.
func Short() string {
	if len(CommitSHA) > 7 {
		return fmt.Sprintf("%s (%s)", Version, CommitSHA[:7])
	}
	return fmt.Sprintf("%s (%s)", Version, CommitSHA)
}
