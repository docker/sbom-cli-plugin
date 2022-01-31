package internal

const valueNotProvided = "[not provided]"

// all variables here are provided as build-time arguments, with clear default values
var version = valueNotProvided
var gitCommit = valueNotProvided
var gitDescription = valueNotProvided
var buildDate = valueNotProvided

// Version defines the application version details (generally from build information)
type Version struct {
	Version        string `json:"version"`        // application semantic version
	GitCommit      string `json:"gitCommit"`      // git SHA at build-time
	GitDescription string `json:"gitDescription"` // output of 'git describe --dirty --always --tags'
	BuildDate      string `json:"buildDate"`      // date of the build
}

// FromBuild provides all version details
func FromBuild() Version {
	return Version{
		Version:        version,
		GitCommit:      gitCommit,
		GitDescription: gitDescription,
		BuildDate:      buildDate,
	}
}
