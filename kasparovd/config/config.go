package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/c4ei/c4exapi/config"
	"github.com/c4ei/c4exapi/version"
	"github.com/c4ei/c4exd/util"
	"github.com/jessevdk/go-flags"
)

const (
	logFilename    = "kasparovd.log"
	errLogFilename = "kasparovd_err.log"
)

var (
	// Default configuration options
	defaultLogDir     = util.AppDataDir("kasparovd", false)
	defaultHTTPListen = "0.0.0.0:8080"
	activeConfig      *Config
)

// ActiveConfig returns the active configuration struct
func ActiveConfig() *Config {
	return activeConfig
}

// Config defines the configuration options for the API server.
type Config struct {
	HTTPListen string `long:"listen" description:"HTTP address to listen on (default: 0.0.0.0:8080)"`
	config.KasparovFlags
}

// Parse parses the CLI arguments and returns a config struct.
func Parse() error {
	activeConfig = &Config{
		HTTPListen: defaultHTTPListen,
	}
	parser := flags.NewParser(activeConfig, flags.HelpFlag)

	_, err := parser.Parse()

	// Show the version and exit if the version flag was specified.
	if activeConfig.ShowVersion {
		appName := filepath.Base(os.Args[0])
		appName = strings.TrimSuffix(appName, filepath.Ext(appName))
		fmt.Println(appName, "version", version.Version())
		os.Exit(0)
	}

	if err != nil {
		return err
	}

	err = activeConfig.ResolveKasparovFlags(parser, defaultLogDir, logFilename, errLogFilename, false)
	if err != nil {
		return err
	}

	return nil
}
