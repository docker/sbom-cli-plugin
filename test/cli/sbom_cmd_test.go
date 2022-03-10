package cli

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/anchore/docker-sbom-cli-plugin/internal"
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
				assertInOutput("an image argument is required"),                        // specific error that should be shown
				assertInOutput("Generate a packaged-based Software Bill Of Materials"), // excerpt from help description
				assertFailingReturnCode,
			},
		},
		{
			name: "json-output-flag",
			args: []string{"sbom", "-o", "json", coverageImage},
			assertions: []traitAssertion{
				assertJsonReport,
				assertJsonDescriptor(internal.SyftName, "v0.41.1"),
				assertNotInOutput("not provided"),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "multiple-output-flags",
			args: []string{"sbom", "-o", "table", "-o", "json=" + tmp + ".tmp/multiple-output-flag-test.json", coverageImage},
			assertions: []traitAssertion{
				assertTableReport,
				assertFileExists(tmp + ".tmp/multiple-output-flag-test.json"),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "table-output-flag",
			args: []string{"sbom", "-o", "table", coverageImage},
			assertions: []traitAssertion{
				assertTableReport,
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "default-output-flag",
			args: []string{"sbom", coverageImage},
			assertions: []traitAssertion{
				assertTableReport,
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "squashed-scope-flag",
			args: []string{"sbom", "-o", "json", "-s", "squashed", coverageImage},
			assertions: []traitAssertion{
				assertPackageCount(20),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "all-layers-scope-flag",
			args: []string{"sbom", "-o", "json", "-s", "all-layers", coverageImage},
			assertions: []traitAssertion{
				assertPackageCount(22),
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "platform-option-wired-up",
			args: []string{"sbom", "--platform", "arm64", "-o", "json", "busybox:1.31"},
			assertions: []traitAssertion{
				assertInOutput("sha256:dcd4bbdd7ea2360002c684968429a2105997c3ce5821e84bfc2703c5ec984971"), // linux/arm64 image digest
				assertSuccessfulReturnCode,
			},
		},
		{
			name: "json-file-flag",
			args: []string{"sbom", "-o", "json", "--file", filepath.Join(tmp, "output-1.json"), coverageImage},
			assertions: []traitAssertion{
				assertSuccessfulReturnCode,
				assertFileOutput(t, filepath.Join(tmp, "output-1.json"),
					assertJsonReport,
				),
			},
		},
		{
			name: "json-output-flag-to-file",
			args: []string{"sbom", "-o", fmt.Sprintf("json=%s", filepath.Join(tmp, "output-2.json")), coverageImage},
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
