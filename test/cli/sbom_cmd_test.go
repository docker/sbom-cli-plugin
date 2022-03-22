package cli

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/sbom-cli-plugin/internal"
)

func TestSBOMCmdFlags(t *testing.T) {
	coverageImage := getFixtureImage(t, "image-pkg-coverage")
	tmp := t.TempDir() + "/"

	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		assertions []traitAssertion
	}{
		{
			name: "no-args-shows-help",
			args: []string{"sbom"},
			assertions: []traitAssertion{
				assertInOutput("an image argument is required"),                                          // specific error that should be shown
				assertInOutput("View the packaged-based Software Bill Of Materials (SBOM) for an image"), // excerpt from help description
				assertFailingReturnCode,
			},
		},
		{
			name: "use-version-option",
			args: []string{"sbom", "version"},
			assertions: []traitAssertion{
				assertInOutput("Application:"),
				assertInOutput("docker-sbom ("),
				assertInOutput("Provider:"),
				assertInOutput("GitDescription:"),
				assertInOutput("syft (v0.42.2)"),
				assertNotInOutput("not provided"),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "use-short-version-option",
			args: []string{"sbom", "--version"},
			assertions: []traitAssertion{
				assertInOutput("sbom-cli-plugin"),
				assertInOutput(", build"),
				assertNotInOutput("not provided"),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "json-format-flag",
			args: []string{"sbom", "--format", "json", coverageImage},
			assertions: []traitAssertion{
				assertJsonReport,
				assertJsonDescriptor(internal.SyftName, "v0.42.2"),
				assertNotInOutput("not provided"),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "table-format-flag",
			args: []string{"sbom", "--format", "table", coverageImage},
			assertions: []traitAssertion{
				assertTableReport,
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "default-format-flag",
			args: []string{"sbom", coverageImage},
			assertions: []traitAssertion{
				assertTableReport,
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "squashed-scope-flag",
			args: []string{"sbom", "--format", "json", "--layers", "squashed", coverageImage},
			assertions: []traitAssertion{
				assertPackageCount(20),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "all-layers-scope-flag",
			args: []string{"sbom", "--format", "json", "--layers", "all-layers", coverageImage},
			assertions: []traitAssertion{
				assertPackageCount(22),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "platform-option-wired-up",
			args: []string{"sbom", "--platform", "arm64", "--format", "json", "busybox:1.31"},
			assertions: []traitAssertion{
				assertInOutput("sha256:dcd4bbdd7ea2360002c684968429a2105997c3ce5821e84bfc2703c5ec984971"), // linux/arm64 image digest
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "json-output-flag",
			args: []string{"sbom", "--format", "json", "--output", filepath.Join(tmp, "output-1.json"), coverageImage},
			assertions: []traitAssertion{
				assertSuccessfulReturnCode,
				assertFileOutput(t, filepath.Join(tmp, "output-1.json"),
					assertJsonReport,
				),
			},
		},
		{
			name: "json-short-output-flag",
			args: []string{"sbom", "--format", "json", "-o", filepath.Join(tmp, "output-2.json"), coverageImage},
			assertions: []traitAssertion{
				assertSuccessfulReturnCode,
				assertFileOutput(t, filepath.Join(tmp, "output-2.json"),
					assertJsonReport,
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, stdout, stderr := runSyft(t, tt.env, tt.args...)
			for _, traitFn := range tt.assertions {
				traitFn(t, stdout, stderr, cmd.ProcessState.ExitCode())
			}
			if t.Failed() {
				t.Log("STDOUT:\n", stdout)
				t.Log("STDERR:\n", stderr)
				t.Log("COMMAND:", strings.Join(cmd.Args, " "))
			}
		})
	}
}
