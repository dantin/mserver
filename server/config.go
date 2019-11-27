package server

import (
	"flag"
	"fmt"
	"os"
)

// Config keeps the configuration of server.
type Config struct {
	*flag.FlagSet

	ListenAddr string `toml:"addr"`

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

	return cfg
}
