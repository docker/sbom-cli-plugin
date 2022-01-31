package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/anchore/docker-sbom-cli-plugin/internal"
	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

const syftBinName = "syft"

func Execute() {
	plugin.Run(
		cmd,
		manager.Metadata{
			SchemaVersion: internal.SchemaVersion,
			Vendor:        "Anchore Inc.",
			Version:       internal.FromBuild().Version,
		},
	)
}

func cmd(_ command.Cli) *cobra.Command {
	return &cobra.Command{
		Use:                "sbom [image]",
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE:               run,
	}
}

func run(cmd *cobra.Command, args []string) error {
	path, err := exec.LookPath(syftBinName)
	if err != nil {
		return fmt.Errorf("%s is not installed", syftBinName)
	}

	// TODO: add limitations such that only the docker daemon can be referenced
	child := exec.Command(path, append([]string{"-v"}, args...)...)
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	return child.Run()
}
