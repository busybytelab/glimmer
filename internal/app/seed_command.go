package app

import (
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"

	"github.com/busybytelab.com/glimmer/data"
	"github.com/busybytelab.com/glimmer/internal/seed"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
)

// seedCommand holds the dependencies and state for seed-related commands
type seedCommand struct {
	pb *pocketbase.PocketBase
}

// handleSeed handles the main seed command execution
func (s *seedCommand) handleSeed(cmd *cobra.Command, args []string) {
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
	if err := seed.RunSeedFromYAML(s.pb, configPath); err != nil {
		log.Fatal().Err(err).Msg("Failed to seed database")
	}

	log.Info().Msg("Seeding completed successfully!")
}

// handlePasswordHash handles the password-hash subcommand execution
func (s *seedCommand) handlePasswordHash(cmd *cobra.Command, args []string) {
	password := args[0]

	// Generate bcrypt hash with cost 10 (same as PocketBase)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate password hash")
	}

	// Generate a tokenKey in PocketBase's format
	// From the sample file, it appears to be a fixed length string of alphanumeric characters
	tokenKey := make([]byte, 32)
	if _, err := crand.Read(tokenKey); err != nil {
		log.Fatal().Err(err).Msg("Failed to generate token key")
	}
	tokenKeyStr := base64.RawURLEncoding.EncodeToString(tokenKey)[:32] // Ensure fixed length

	// Output the hash in PocketBase's format: bcrypt_hash|tokenKey
	fmt.Printf("%s|%s\n", string(hash), tokenKeyStr)
}

// setupSeedCommand configures the seed command for the application
func setupSeedCommand(pb *pocketbase.PocketBase) {
	log.Trace().Msg("Setting up seed command...")

	// Create command handler
	handler := &seedCommand{pb: pb}

	// Create the seed command
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed the database with test data from a YAML configuration file",
		Run:   handler.handleSeed,
	}

	// Add flags to the command
	seedCmd.Flags().String("config", "", "Path to seed YAML config file (default: uses embedded sample-v0.1.yaml)")

	// Create the password-hash subcommand
	passwordHashCmd := &cobra.Command{
		Use:   "password-hash [password]",
		Short: "Generate a bcrypt password hash for use in seed files",
		Args:  cobra.ExactArgs(1),
		Run:   handler.handlePasswordHash,
	}

	// Add the password-hash subcommand to the seed command
	seedCmd.AddCommand(passwordHashCmd)

	// Add the seed command to PocketBase's root command
	pb.RootCmd.AddCommand(seedCmd)

	log.Trace().Msg("Seed command setup completed")
}
