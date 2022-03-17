package cmd

import (
	"fmt"

	"github.com/docker/sbom-cli-plugin/internal"
	"github.com/docker/sbom-cli-plugin/internal/version"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	c := &cobra.Command{
		Use:           "version",
		Short:         "Show Docker sbom version information",
		Args:          cobra.NoArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(_ *cobra.Command, args []string) error {
			report := tprintf(`Application:        {{ .Name }} ({{ .Version.Version }})
Provider:           {{ .SyftName }} ({{ .SyftVersion }})
GitCommit:          {{ .GitCommit }}
GitDescription:     {{ .GitDescription }}
Platform:           {{ .Platform }}
`, struct {
				Name     string
				SyftName string
				version.Version
			}{
				Name:     internal.BinaryName,
				SyftName: internal.SyftName,
				Version:  version.FromBuild(),
			})

			fmt.Print(report)
			return nil
		},
	}

	return c
}
