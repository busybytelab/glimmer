package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/busybytelab.com/glimmer/internal/seed"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "", "Path to seed YAML config file (default: test/data/seed_data.yaml)")
	flag.Parse()

	// Log the start of the seeding process
	log.Println("Starting YAML-based database seeding...")

	// Run the seed from YAML
	seed.RunSeedFromYAML(*configPath)

	fmt.Println("Seeding completed successfully!")
	os.Exit(0)
}
