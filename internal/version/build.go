/*
Package version contains all build time metadata (version, build time, git commit, etc).
*/
package version

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/anchore/docker-sbom-cli-plugin/internal/log"
)

const valueNotProvided = "[not provided]"

// all variables here are provided as build-time arguments, with clear default values
var (
	version        = valueNotProvided
	gitCommit      = valueNotProvided
	gitDescription = valueNotProvided
	buildDate      = valueNotProvided
)

// Version defines the application version details (generally from build information)
type Version struct {
	Version        string `json:"version"`        // application semantic version
	SyftVersion    string `json:"syftVersion"`    // the version of syft being used by the docker-sbom-cli-plugin
	GitCommit      string `json:"gitCommit"`      // git SHA at build-time
	GitDescription string `json:"gitDescription"` // output of 'git describe --dirty --always --tags'
	BuildDate      string `json:"buildDate"`      // date of the build
	GoVersion      string `json:"goVersion"`      // go runtime version at build-time
	Compiler       string `json:"compiler"`       // compiler used at build-time
	Platform       string `json:"platform"`       // GOOS and GOARCH at build-time
}

// FromBuild provides all version details
func FromBuild() Version {
	return Version{
		Version:        version,
		SyftVersion:    syftVersion(),
		GitCommit:      gitCommit,
		GitDescription: gitDescription,
		BuildDate:      buildDate,
		GoVersion:      runtime.Version(),
		Compiler:       runtime.Compiler,
		Platform:       fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func syftVersion() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Warn("unable to find the buildinfo section of the binary (syft version is unknown)")
		return valueNotProvided
	}

	for _, d := range buildInfo.Deps {
		if d.Path == "github.com/anchore/syft" {
			return d.Version
		}
	}

	log.Warn("unable to find 'github.com/anchore/syft' from the buildinfo section of the binary")

	return valueNotProvided
}
