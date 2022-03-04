package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/docker/cli/cli-plugins/manager"
	"github.com/docker/cli/cli-plugins/plugin"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/wagoodman/go-partybus"

	"github.com/anchore/docker-sbom-cli-plugin/internal"
	"github.com/anchore/docker-sbom-cli-plugin/internal/bus"
	"github.com/anchore/docker-sbom-cli-plugin/internal/config"
	"github.com/anchore/docker-sbom-cli-plugin/internal/log"
	"github.com/anchore/docker-sbom-cli-plugin/internal/logger"
	"github.com/anchore/docker-sbom-cli-plugin/internal/version"
	"github.com/anchore/stereoscope"
	"github.com/anchore/syft/syft"
)

var (
	appConfig         *config.Application
	eventBus          *partybus.Bus
	eventSubscription *partybus.Subscription
)

func init() {
	cobra.OnInitialize(
		initAppConfig,
		initLogging,
		logAppConfig,
		logAppVersion,
		initEventBus,
	)
}

func Execute() {
	plugin.Run(
		cmd,
		manager.Metadata{
			SchemaVersion: internal.SchemaVersion,
			Vendor:        "Anchore Inc.",
			Version:       version.FromBuild().Version,
		},
	)
}

func initAppConfig() {
	cfg, err := config.LoadApplicationConfig(viper.GetViper(), cliOnlyOpts)
	if err != nil {
		fmt.Printf("failed to load application config: \n\t%+v\n", err)
		os.Exit(1)
	}

	appConfig = cfg
}

func initLogging() {
	cfg := logger.LogrusConfig{
		EnableConsole: (appConfig.Log.FileLocation == "" || appConfig.CliOptions.Verbosity > 0) && !appConfig.Quiet,
		EnableFile:    appConfig.Log.FileLocation != "",
		Level:         appConfig.Log.LevelOpt,
		Structured:    appConfig.Log.Structured,
		FileLocation:  appConfig.Log.FileLocation,
	}

	logWrapper := logger.NewLogrusLogger(cfg)
	syft.SetLogger(logWrapper)
	stereoscope.SetLogger(&logger.LogrusNestedLogger{
		Logger: logWrapper.Logger.WithField("from-lib", "stereoscope"),
	})
	log.Log = logWrapper
}

func logAppConfig() {
	log.Debugf("application config:\n%+v", color.Magenta.Sprint(appConfig.String()))
}

func initEventBus() {
	eventBus = partybus.NewBus()
	eventSubscription = eventBus.Subscribe()

	stereoscope.SetBus(eventBus)
	syft.SetBus(eventBus)
	bus.SetPublisher(eventBus)
}

func logAppVersion() {
	versionInfo := version.FromBuild()
	log.Infof("%s version: %s", internal.SyftName, versionInfo.SyftVersion)

	var fields map[string]interface{}
	bytes, err := json.Marshal(versionInfo)
	if err != nil {
		return
	}
	err = json.Unmarshal(bytes, &fields)
	if err != nil {
		return
	}

	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for idx, field := range keys {
		value := fields[field]
		branch := "├──"
		if idx == len(fields)-1 {
			branch = "└──"
		}
		log.Debugf("  %s %s: %s", branch, field, value)
	}
}
