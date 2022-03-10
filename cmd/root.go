package cmd

import (
	"errors"
	"fmt"

	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wagoodman/go-partybus"

	"github.com/anchore/docker-sbom-cli-plugin/internal"
	"github.com/anchore/docker-sbom-cli-plugin/internal/bus"
	"github.com/anchore/docker-sbom-cli-plugin/internal/config"
	"github.com/anchore/docker-sbom-cli-plugin/internal/log"
	"github.com/anchore/docker-sbom-cli-plugin/internal/ui"
	"github.com/anchore/docker-sbom-cli-plugin/internal/version"
	"github.com/anchore/stereoscope"
	"github.com/anchore/stereoscope/pkg/image"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/event"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
)

const helpExample = `docker sbom [options] [IMAGE]

Examples:
  docker sbom alpine:latest                                        a summary of discovered packages
  docker sbom alpine:latest -o syft-json                           show all possible cataloging details
  docker sbom alpine:latest -o syft-json --file sbom.json          write report output to a file
  docker sbom alpine:latest -o table -o sbom.json=cyclonedx-json   report the SBOM in multiple formats
  docker sbom alpine:latest --exclude /lib  --exclude '**/*.db'    ignore one or more paths in the image
  docker sbom alpine:latest -v                                     show logging output
  docker sbom alpine:latest -vv                                    show verbose debug logs`

var cliOnlyOpts = config.CliOnlyOptions{}

func cmd(_ command.Cli) *cobra.Command {
	c := &cobra.Command{
		Use:               "sbom",
		Short:             "Generate a package SBOM",
		Long:              "Generate a packaged-based Software Bill Of Materials (SBOM) from Docker images",
		Example:           helpExample,
		Args:              validateInputArgs,
		SilenceUsage:      true,
		SilenceErrors:     true,
		RunE:              run,
		ValidArgsFunction: dockerImageValidArgsFunction,
	}

	// setting the version template to just print out the string since we already have a templatized version string
	c.SetVersionTemplate(fmt.Sprintf("%s {{.Version}}\n", internal.ApplicationName))

	setPackageFlags(c.Flags())

	if err := bindConfigOptions(c.Flags()); err != nil {
		panic(fmt.Errorf("unable to bind config options: %w", err))
	}

	return c
}

func setPackageFlags(flags *pflag.FlagSet) {
	// Universal options ///////////////////////////////////////////////////////

	flags.StringVarP(&cliOnlyOpts.ConfigPath, "config", "c", "", "application config file")

	flags.BoolP(
		"quiet", "q", false,
		"suppress all logging output",
	)

	flags.CountVarP(&cliOnlyOpts.Verbosity, "verbose", "v", "increase verbosity (-v = info, -vv = debug)")

	// Formatting & Input options //////////////////////////////////////////////
	flags.StringP(
		"scope", "s", cataloger.DefaultSearchConfig().Scope.String(),
		fmt.Sprintf("[experimental] selection of layers to catalog, options=%v", source.AllScopes))

	flags.StringArrayP(
		"output", "o", formatAliases(syft.TableFormatID),
		fmt.Sprintf("report output format, options=%v", formatAliases(syft.FormatIDs()...)),
	)

	flags.StringP(
		"file", "", "",
		"file to write the default report output to (default is STDOUT)",
	)

	flags.StringArrayP(
		"exclude", "", nil,
		"exclude paths from being scanned using a glob expression",
	)

	flags.StringP(
		"platform", "", "",
		"an optional platform specifier for container image sources (e.g. 'linux/arm64', 'linux/arm64/v8', 'arm64', 'linux')",
	)
}

func bindConfigOptions(flags *pflag.FlagSet) error {
	// Universal options ///////////////////////////////////////////////////////

	if err := viper.BindPFlag("quiet", flags.Lookup("quiet")); err != nil {
		return err
	}

	// Formatting & Input options //////////////////////////////////////////////

	if err := viper.BindPFlag("output", flags.Lookup("output")); err != nil {
		return err
	}

	if err := viper.BindPFlag("package.cataloger.scope", flags.Lookup("scope")); err != nil {
		return err
	}

	if err := viper.BindPFlag("file", flags.Lookup("file")); err != nil {
		return err
	}

	if err := viper.BindPFlag("exclude", flags.Lookup("exclude")); err != nil {
		return err
	}

	if err := viper.BindPFlag("platform", flags.Lookup("platform")); err != nil {
		return err
	}

	return nil
}
func validateInputArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// in the case that no arguments are given we want to show the help text and return with a non-0 return code.
		if err := cmd.Help(); err != nil {
			return fmt.Errorf("unable to display help: %w", err)
		}
		return fmt.Errorf("an image argument is required")
	}

	return cobra.ExactArgs(1)(cmd, args)
}

func run(_ *cobra.Command, args []string) error {
	writer, err := makeWriter(appConfig.Output, appConfig.File)
	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Warnf("unable to write to report destination: %+v", err)
		}
	}()

	si := source.Input{
		UserInput:   args[0],
		Scheme:      source.ImageScheme,
		ImageSource: image.DockerDaemonSource,
		Location:    args[0],
		Platform:    appConfig.Platform,
	}

	return eventLoop(
		sbomExecWorker(si, writer),
		setupSignals(),
		eventSubscription,
		stereoscope.Cleanup,
		ui.Select(isVerbose(), appConfig.Quiet)...,
	)
}

func isVerbose() (result bool) {
	isPipedInput, err := internal.IsPipedInput()
	if err != nil {
		// since we can't tell if there was piped input we assume that there could be to disable the ETUI
		log.Warnf("unable to determine if there is piped input: %+v", err)
		return true
	}
	// verbosity should consider if there is piped input (in which case we should not show the ETUI)
	return appConfig.CliOptions.Verbosity > 0 || isPipedInput
}

func generateSBOM(src *source.Source) (*sbom.SBOM, error) {
	s := sbom.SBOM{
		Source: src.Metadata,
		Descriptor: sbom.Descriptor{
			Name:          internal.SyftName,
			Version:       version.FromBuild().SyftVersion,
			Configuration: appConfig,
		},
	}

	packageCatalog, relationships, theDistro, err := syft.CatalogPackages(src, appConfig.Package.ToConfig())
	if err != nil {
		return nil, fmt.Errorf("unable to catalog packages: %w", err)
	}

	s.Artifacts.PackageCatalog = packageCatalog
	s.Artifacts.LinuxDistribution = theDistro
	s.Relationships = relationships

	return &s, nil
}

func sbomExecWorker(si source.Input, writer sbom.Writer) <-chan error {
	errs := make(chan error)
	go func() {
		defer close(errs)

		src, cleanup, err := source.New(si, nil, appConfig.Exclusions)
		if cleanup != nil {
			defer cleanup()
		}
		if err != nil {
			errs <- fmt.Errorf("failed to construct source from user input %q: %w", si.UserInput, err)
			return
		}

		s, err := generateSBOM(src)
		if err != nil {
			errs <- err
			return
		}

		if err != nil {
			errs <- errors.New("could not produce an sbom")
			return
		}

		bus.Publish(partybus.Event{
			Type:  event.Exit,
			Value: func() error { return writer.Write(*s) },
		})
	}()
	return errs
}
