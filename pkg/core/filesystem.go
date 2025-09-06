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

	// Create default TOC template
	tocTemplateFile := filepath.Join(tplDir, "999_inhaltsverzeichnis_template.abc")
	tocTemplateContent := `X:1
T:Inhaltsverzeichnis
M:4/4
L:1/4
K:C
W:{{TOC}}
`
	
	if err := os.WriteFile(tocTemplateFile, []byte(tocTemplateContent), 0644); err != nil {
		return fmt.Errorf("failed to create TOC template for %s: %w", shortName, err)
	}

	return nil
}
