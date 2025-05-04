package main

import (
	"embed"
	"io/fs"
)

//go:embed all:ui/dist
var distDir embed.FS

// DistDirFS contains the embedded dist directory files (without the "dist" prefix)
var DistDirFS, _ = fs.Sub(distDir, "ui/dist")
