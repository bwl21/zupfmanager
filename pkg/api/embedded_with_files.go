//go:build embed_frontend

package api

import (
	"embed"
	"io/fs"
)

// Embedded frontend files - only included when building with embed_frontend tag
//go:embed frontend/dist
var frontendFS embed.FS

// Embedded configuration files
//go:embed default-project-config.json
var defaultConfigFS embed.FS

// GetFrontendFS returns the embedded frontend filesystem
func GetFrontendFS() (fs.FS, error) {
	// Return the subdirectory containing the actual frontend files
	return fs.Sub(frontendFS, "frontend/dist")
}

// GetDefaultConfigFS returns the embedded default configuration filesystem
func GetDefaultConfigFS() (fs.FS, error) {
	return defaultConfigFS, nil
}
