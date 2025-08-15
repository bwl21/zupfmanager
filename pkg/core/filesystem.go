package core

import (
	"fmt"
	"os"
	"path/filepath"
)

// fileSystemService implements FileSystemService interface
type fileSystemService struct{}

// NewFileSystemService creates a new filesystem service
func NewFileSystemService() FileSystemService {
	return &fileSystemService{}
}

// CreateProjectDirectories creates the required directory structure for a project
func (f *fileSystemService) CreateProjectDirectories(shortName string) error {
	if shortName == "" {
		return fmt.Errorf("short name cannot be empty")
	}

	projectDir := shortName
	tplDir := filepath.Join(projectDir, "tpl")
	
	if err := os.MkdirAll(tplDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directories for %s: %w", shortName, err)
	}

	return nil
}
