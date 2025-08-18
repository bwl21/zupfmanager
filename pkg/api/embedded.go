//go:build !embed_frontend

package api

import (
	"errors"
	"io/fs"
)

// GetFrontendFS returns an error for non-embedded builds
func GetFrontendFS() (fs.FS, error) {
	return nil, errors.New("frontend not embedded in this build")
}
