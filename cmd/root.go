package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/docker/cli/cli"
	"github.com/docker/cli/cli/command"
	"github.com/docker/sbom-cli-plugin/internal"
	"github.com/docker/sbom-cli-plugin/internal/bus"
	"github.com/docker/sbom-cli-plugin/internal/log"
	"github.com/docker/sbom-cli-plugin/internal/ui"
	"github.com/docker/sbom-cli-plugin/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/wagoodman/go-partybus"

	"github.com/anchore/stereoscope"
	"github.com/anchore/stereoscope/pkg/file"
	"github.com/anchore/stereoscope/pkg/image"
	stereoscopeDocker "github.com/anchore/stereoscope/pkg/image/docker"
	"github.com/anchore/syft/syft"
	"github.com/anchore/syft/syft/event"
	"github.com/anchore/syft/syft/pkg/cataloger"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"
)

const (
	helpExample = `
  docker sbom alpine:latest                                          a summary of discovered packages
  docker sbom alpine:latest --format syft-json                       show all possible cataloging details
  docker sbom alpine:latest --output sbom.txt                        write report output to a file
  docker sbom alpine:latest --exclude /lib  --exclude '**/*.db'      ignore one or more paths/globs in the image
`
	shortDescription = "View the packaged-based Software Bill Of Materials (SBOM) for an image"
)

func cmd(dockerCli command.Cli) *cobra.Command {
	c := &cobra.Command{
		Use:           "sbom",
		Short:         shortDescription,
		Long:          shortDescription + ".\n\nEXPERIMENTAL: The flags and outputs of this command may change. Leave feedback on https://github.com/docker/sbom-cli-plugin.",
		Example:       helpExample,
		Args:          validateInputArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.FromBuild().Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return cli.StatusError{
				StatusCode: 1,
				Status:     `error: docker sbom has been removed, please use "docker scout sbom" command instead`,
			}
		},
	}

	c.SetVersionTemplate(fmt.Sprintf("%s {{.Version}}, build %s\n", internal.ApplicationName, version.FromBuild().GitCommit))

	setPackageFlags(c.Flags())

	if err := bindConfigOptions(c.Flags()); err != nil {
		panic(fmt.Errorf("unable to bind config options: %w", err))
	}

	c.AddCommand(versionCmd())

	return c
}

func tprintf(tmpl string, data interface{}) string {
	t := template.Must(template.New("").Parse(tmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return ""
	}
	return buf.String()
}

func allScopes() (result []string) {
	for _, s := range source.AllScopes {
		result = append(result, cleanScope(s))
	}
	return result
}

func cleanScope(s source.Scope) string {
	var opt string
	switch s {
	case source.AllLayersScope:
		opt = "all"
	case source.SquashedScope:
		opt = "squashed"
	default:
		opt = strings.ToLower(string(s))
	}
	return opt
}

func setPackageFlags(flags *pflag.FlagSet) {
	flags.BoolP(
		"quiet", "q", false,
		"suppress all non-report output",
	)

	flags.StringP(
		"layers", "", cleanScope(cataloger.DefaultSearchConfig().Scope),
		fmt.Sprintf("[experimental] selection of layers to catalog, options=%v", allScopes()),
	)

	flags.StringP(
		"format", "", formatAliases(syft.TableFormatID)[0],
		fmt.Sprintf("report output format, options=%v", formatAliases(syft.FormatIDs()...)),
	)

	flags.StringP(
		"output", "o", "",
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

	flags.BoolP(
		"debug", "D", false,
		"show debug logging",
	)
}

func bindConfigOptions(flags *pflag.FlagSet) error {
	if err := viper.BindPFlag("quiet", flags.Lookup("quiet")); err != nil {
		return err
	}

	if err := viper.BindPFlag("output", flags.Lookup("output")); err != nil {
		return err
	}

	if err := viper.BindPFlag("package.cataloger.scope", flags.Lookup("layers")); err != nil {
		return err
	}

	if err := viper.BindPFlag("format", flags.Lookup("format")); err != nil {
		return err
	}

	if err := viper.BindPFlag("exclude", flags.Lookup("exclude")); err != nil {
		return err
	}

	if err := viper.BindPFlag("platform", flags.Lookup("platform")); err != nil {
		return err
	}

	if err := viper.BindPFlag("debug", flags.Lookup("debug")); err != nil {
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

type runner struct {
	client command.Cli
}

func newRunner(client command.Cli) runner {
	return runner{
		client: client,
	}
}

func (r runner) run(_ *cobra.Command, args []string) error {
	writer, err := makeWriter([]string{appConfig.Format}, appConfig.Output)
	if err != nil {
		return err
	}

	defer func() {
		if err := writer.Close(); err != nil {
			log.Warnf("unable to write to report destination: %+v", err)
		}
	}()

	var platform *image.Platform
	if appConfig.Platform != "" {
		platform, err = image.NewPlatform(appConfig.Platform)
		if err != nil {
			return fmt.Errorf("invalid platform provided: %w", err)
		}
	}

	cleanImageName, err := cleanImageReference(args[0])
	if err != nil {
		return nil
	}

	return eventLoop(
		sbomExecWorker(cleanImageName, r.client, platform, writer),
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
	return appConfig.Debug || isPipedInput
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

func sbomExecWorker(imageName string, dockerCli command.Cli, platform *image.Platform, writer sbom.Writer) <-chan error {
	errs := make(chan error)
	go func() {
		defer close(errs)

		tempGen := file.NewTempDirGenerator(internal.ApplicationName)

		provider := stereoscopeDocker.NewProviderFromDaemon(
			imageName,
			tempGen,
			dockerCli.Client(),
			platform,
		)
		img, err := provider.Provide(context.Background())
		defer func() {
			if err := tempGen.Cleanup(); err != nil {
				log.Warnf("failed to clean up image: %+v", err)
			}
		}()
		if err != nil {
			errs <- fmt.Errorf("failed to fetch the image %q: %w", imageName, err)
			return
		}

		err = img.Read()
		if err != nil {
			errs <- fmt.Errorf("failed to read the image %q: %w", imageName, err)
			return
		}

		src, err := source.NewFromImage(img, imageName)
		if err != nil {
			errs <- fmt.Errorf("failed to construct source from user input %q: %w", imageName, err)
			return
		}
		src.Exclusions = appConfig.Exclusions

		s, err := generateSBOM(&src)
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
