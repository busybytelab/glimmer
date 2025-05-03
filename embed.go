package main

import (
	"embed"
	"io/fs"
)

//go:embed pb_public
var publicFS embed.FS

// PublicFS returns the embedded filesystem with static web files.
// This function provides the embedded pb_public directory for serving static content.
func PublicFS() fs.FS {
	// Extract the embedded subdirectory to use as the root
	fsys, err := fs.Sub(publicFS, "pb_public")
	if err != nil {
		// This should never happen, but handle it just in case
		panic("failed to extract embedded pb_public directory: " + err.Error())
	}
	return fsys
}
