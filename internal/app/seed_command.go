package app

import (
	"io/fs"
	"os"

	"github.com/busybytelab.com/glimmer/data"
	"github.com/busybytelab.com/glimmer/internal/seed"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// setupSeedCommand configures the seed command for the application
func setupSeedCommand(pb *pocketbase.PocketBase) {
	log.Trace().Msg("Setting up seed command...")

	// Create the seed command
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with test data from a YAML configuration file",
		Run: func(cmd *cobra.Command, args []string) {
			// Parse command line flags
			configPath := cmd.Flag("config").Value.String()
			if configPath == "" {
				// Use embedded data if no config path is provided
				configPath = "sample-v0.1.yaml"
				// Create a temporary file to store the embedded YAML
				tmpFile, err := os.CreateTemp("", "seed-*.yaml")
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to create temporary file")
				}
				defer os.Remove(tmpFile.Name())

				// Read the embedded YAML file
				data, err := fs.ReadFile(data.SeedDirFS, configPath)
				if err != nil {
					log.Fatal().Err(err).Msg("Failed to read embedded seed data")
				}

				// Write the embedded YAML to the temporary file
				if _, err := tmpFile.Write(data); err != nil {
					log.Fatal().Err(err).Msg("Failed to write seed data to temporary file")
				}
				tmpFile.Close()

				configPath = tmpFile.Name()
			}

			// Log the start of the seeding process
			log.Info().Msg("Starting YAML-based database seeding...")

			// Run the seed from YAML
			seed.RunSeedFromYAML(configPath)

			log.Info().Msg("Seeding completed successfully!")
		},
	}

	// Add flags to the command
	seedCmd.Flags().String("config", "", "Path to seed YAML config file (default: uses embedded sample-v0.1.yaml)")

	// Add the seed command to PocketBase's root command
	pb.RootCmd.AddCommand(seedCmd)

	log.Trace().Msg("Seed command setup completed")
}
