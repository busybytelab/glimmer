package main

import (
	"github.com/rs/zerolog/log"

	"github.com/busybytelab.com/glimmer/internal/app"
	_ "github.com/busybytelab.com/glimmer/internal/migrations" // register migrations
)

// Version holds the application version.
// This is set during build using ldflags.
var Version = "dev"

func main() {
	application := app.New(DistDirFS)

	if err := application.Initialize(); err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize glimmer application")
	}

	if err := application.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start glimmer application")
	} else {
		log.Info().Str("version", Version).Msg("Glimmer application started successfully")
	}
}
