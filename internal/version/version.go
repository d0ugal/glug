package version

import (
	"fmt"
	"runtime"
)

// These variables are set during build time using ldflags
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

// Info holds version information
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
}

// Get returns the version information
func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("glug %s (commit: %s, built: %s, go: %s)",
		i.Version, i.Commit, i.BuildDate, i.GoVersion)
}
