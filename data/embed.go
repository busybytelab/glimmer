package data

import (
	"embed"
	"io/fs"
)

//go:embed all:seed
var seedDir embed.FS

// SeedDirFS contains the embedded seed directory files (without the "seed" prefix)
var SeedDirFS, _ = fs.Sub(seedDir, "seed")
