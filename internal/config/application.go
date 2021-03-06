package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type defaultValueLoader interface {
	loadDefaultValues(*viper.Viper)
}

type parser interface {
	parseConfigValues() error
}

// Application is the main syft application configuration.
type Application struct {
	Package    pkg      `yaml:"package" json:"package" mapstructure:"package"`    // package cataloging related options
	Exclusions []string `yaml:"exclude" json:"exclude" mapstructure:"exclude"`    // --exclude, ignore paths within an image
	Platform   string   `yaml:"platform" json:"platform" mapstructure:"platform"` // --platform, override OS and architecture from image
	Output     string   `yaml:"output" json:"output" mapstructure:"output"`       // --output, the file to write report output to
	Format     string   `yaml:"format" json:"format" mapstructure:"format"`       // --format, the format to use for output
	Quiet      bool     `yaml:"quiet" json:"quiet" mapstructure:"quiet"`          // -q, indicates to not show any status output to stderr (ETUI or logging UI)
	Log        logging  `yaml:"log" json:"log" mapstructure:"log"`                // all logging-related options
	Debug      bool     `yaml:"debug" json:"debug" mapstructure:"debug"`          // -D/--debug, enable debug logging
}

func newApplicationConfig(v *viper.Viper) *Application {
	config := &Application{}
	config.loadDefaultValues(v)
	return config
}

// LoadApplicationConfig populates the given viper object with a default application config values
func LoadApplicationConfig(v *viper.Viper) (*Application, error) {
	// the user may not have a config, and this is OK, we can use the default config + default cobra cli values instead
	config := newApplicationConfig(v)

	// TODO: in the future when we have a user-modifiable configuration, reading such contents would be here

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	if err := config.parseConfigValues(); err != nil {
		return nil, fmt.Errorf("invalid application config: %w", err)
	}

	return config, nil
}

// init loads the default configuration values into the viper instance (before the config values are read and parsed).
func (cfg Application) loadDefaultValues(v *viper.Viper) {
	// for each field in the configuration struct, see if the field implements the defaultValueLoader interface and invoke it if it does
	value := reflect.ValueOf(cfg)
	for i := 0; i < value.NumField(); i++ {
		// note: the defaultValueLoader method receiver is NOT a pointer receiver.
		if loadable, ok := value.Field(i).Interface().(defaultValueLoader); ok {
			// the field implements defaultValueLoader, call it
			loadable.loadDefaultValues(v)
		}
	}
}

func (cfg *Application) parseConfigValues() error {
	// parse application config options
	for _, optionFn := range []func() error{
		cfg.parseLogLevelOption,
	} {
		if err := optionFn(); err != nil {
			return err
		}
	}

	// parse nested config options
	// for each field in the configuration struct, see if the field implements the parser interface
	// note: the app config is a pointer, so we need to grab the elements explicitly (to traverse the address)
	value := reflect.ValueOf(cfg).Elem()
	for i := 0; i < value.NumField(); i++ {
		// note: since the interface method of parser is a pointer receiver we need to get the value of the field as a pointer.
		if parsable, ok := value.Field(i).Addr().Interface().(parser); ok {
			// the field implements parser, call it
			if err := parsable.parseConfigValues(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cfg *Application) parseLogLevelOption() error {
	switch {
	case cfg.Quiet:
		cfg.Log.LevelOpt = logrus.PanicLevel
	case cfg.Log.Level != "":
		lvl, err := logrus.ParseLevel(strings.ToLower(cfg.Log.Level))
		if err != nil {
			return fmt.Errorf("bad log level configured (%q): %w", cfg.Log.Level, err)
		}

		cfg.Log.LevelOpt = lvl
	case cfg.Debug:
		cfg.Log.LevelOpt = logrus.DebugLevel
	default:
		cfg.Log.LevelOpt = logrus.WarnLevel
	}

	return nil
}

func (cfg Application) String() string {
	// yaml is pretty human friendly (at least when compared to json)
	appCfgStr, err := yaml.Marshal(&cfg)

	if err != nil {
		return err.Error()
	}

	return string(appCfgStr)
}
