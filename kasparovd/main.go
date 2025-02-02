package main

import (
	"fmt"
	"os"

	"github.com/c4ei/c4exd/util/profiling"

	"github.com/pkg/errors"

	"github.com/c4ei/c4exapi/database"
	"github.com/c4ei/c4exapi/jsonrpc"
	"github.com/c4ei/c4exapi/kasparovd/config"
	"github.com/c4ei/c4exapi/kasparovd/server"
	"github.com/c4ei/c4exapi/version"
	"github.com/c4ei/c4exd/signal"
	"github.com/c4ei/c4exd/util/panics"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	defer panics.HandlePanic(log, "main", nil)
	interrupt := signal.InterruptListener()

	err := config.Parse()
	if err != nil {
		errString := fmt.Sprintf("Error parsing command-line arguments: %s\n", err)
		_, fErr := fmt.Fprint(os.Stderr, errString)
		if fErr != nil {
			panic(errString)
		}
		return
	}

	// Show version at startup.
	log.Infof("Version %s", version.Version())

	// Start the profiling server if required
	if config.ActiveConfig().Profile != "" {
		profiling.Start(config.ActiveConfig().Profile, log)
	}

	err = database.Connect(&config.ActiveConfig().KasparovFlags)
	if err != nil {
		panic(errors.Errorf("Error connecting to database: %s", err))
	}
	defer func() {
		err := database.Close()
		if err != nil {
			panic(errors.Errorf("Error closing the database: %s", err))
		}
	}()

	err = jsonrpc.Connect(&config.ActiveConfig().KasparovFlags, false)
	if err != nil {
		panic(errors.Errorf("Error connecting to servers: %s", err))
	}
	defer jsonrpc.Close()

	shutdownServer := server.Start(config.ActiveConfig().HTTPListen)
	defer shutdownServer()

	<-interrupt
}
