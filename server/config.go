package server

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/dantin/mserver/pkg/logutil"
	"github.com/juju/errors"
)

// Config keeps the configuration of server.
type Config struct {
	*flag.FlagSet

	ListenAddr string `toml:"addr"`

	Log logutil.LogConfig `toml:"log"`

	configFile   string
	printVersion bool
}

// NewConfig creates an instance of server configuration.
func NewConfig() *Config {
	cfg := &Config{}

	cfg.FlagSet = flag.NewFlagSet(DefaultName, flag.ContinueOnError)
	fs := cfg.FlagSet
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", DefaultName)
		fs.PrintDefaults()
	}

	// server related configuration.
	fs.StringVar(&cfg.ListenAddr, "addr", defaultListenAddr, "addr(e.g. 'host:port') to listen on the specific port")
	// log related configuration.
	fs.StringVar(&cfg.Log.Level, "L", "info", "log level: debug, info, warn, error, fatal")
	fs.StringVar(&cfg.Log.File.Filename, "log-file", "", "log file path")
	// misc related configuration.
	fs.StringVar(&cfg.configFile, "config", "", fmt.Sprintf("path to the %s's configuration file", DefaultName))
	fs.BoolVar(&cfg.printVersion, "V", false, "print version info")

	return cfg
}

// Parse checks all configuration from command-line flags.
func (cfg *Config) Parse(args []string) error {
	// parse first to get config file.
	err := cfg.FlagSet.Parse(args)
	switch err {
	case nil:
	case flag.ErrHelp:
		os.Exit(0)
	default:
		os.Exit(2)
	}

	if cfg.printVersion {
		printVersionInfo()
		os.Exit(0)
	}

	// load configuration if specified.
	if cfg.configFile != "" {
		if err := cfg.configFromFile(cfg.configFile); err != nil {
			return errors.Trace(err)
		}
	}
	// parse again to replace config with command line options.
	cfg.FlagSet.Parse(args)
	if len(cfg.FlagSet.Args()) > 0 {
		return errors.Errorf("%q is not a valid flag", cfg.FlagSet.Arg(0))
	}

	return nil
}

func (cfg *Config) configFromFile(path string) error {
	_, err := toml.DecodeFile(path, cfg)
	return errors.Trace(err)
}
