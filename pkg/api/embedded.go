package api

import (
	"embed"
	"io/fs"
)

// Embedded frontend files
//go:embed frontend/dist
var frontendFS embed.FS

// GetFrontendFS returns the embedded frontend filesystem
func GetFrontendFS() (fs.FS, error) {
	// Return the subdirectory containing the actual frontend files
	return fs.Sub(frontendFS, "frontend/dist")
}
