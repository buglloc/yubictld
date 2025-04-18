package main

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"go.uber.org/automaxprocs/maxprocs"

	"github.com/buglloc/yubictld/internal/commands"
)

func main() {
	if _, err := maxprocs.Set(); err != nil {
		log.Error().Err(err).Msg("set GOMAXPROCS")
	}

	if err := commands.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
