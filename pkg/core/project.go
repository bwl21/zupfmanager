package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
)

// ProjectService handles project-related operations
type ProjectService struct {
	db *database.Client
}

// NewProjectService creates a new project service
func NewProjectService() (*ProjectService, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}
	return &ProjectService{db: db}, nil
}

// CreateProjectRequest represents the data needed to create a project
type CreateProjectRequest struct {
	Title         string
	ShortName     string
	ConfigFile    string
	DefaultConfig bool
}

// CreateProject creates a new project with the given parameters
func (s *ProjectService) CreateProject(ctx context.Context, req CreateProjectRequest) (*ent.Project, error) {
	var config map[string]interface{}

	if req.ConfigFile != "" {
		configFile, err := os.ReadFile(req.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := json.Unmarshal(configFile, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	} else {
		config = map[string]interface{}{}
		if req.DefaultConfig {
			configFile, err := os.ReadFile("default-project-config.json")
			if err != nil {
				return nil, fmt.Errorf("failed to read default config file: %w", err)
			}
			if err := json.Unmarshal(configFile, &config); err != nil {
				return nil, fmt.Errorf("failed to parse default config JSON: %w", err)
			}
		}
	}

	project, err := s.db.CreateOrUpdateProject(ctx, 0, req.Title, req.ShortName, config)
	if err != nil {
		return nil, err
	}

	// Create directory structure
	projectDir := req.ShortName
	tplDir := projectDir + "/tpl"
	if err := os.MkdirAll(tplDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	return project, nil
}

// ListProjects returns all projects
func (s *ProjectService) ListProjects(ctx context.Context) ([]*ent.Project, error) {
	return s.db.Project.Query().All(ctx)
}

// GetProject returns a project by ID
func (s *ProjectService) GetProject(ctx context.Context, id int) (*ent.Project, error) {
	return s.db.GetProject(ctx, id)
}

// UpdateProjectRequest represents the data needed to update a project
type UpdateProjectRequest struct {
	ID            int
	Title         string
	ShortName     string
	ConfigFile    string
	DefaultConfig bool
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(ctx context.Context, req UpdateProjectRequest) (*ent.Project, error) {
	var config map[string]interface{}

	if req.ConfigFile != "" {
		configFile, err := os.ReadFile(req.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := json.Unmarshal(configFile, &config); err != nil {
			return nil, fmt.Errorf("failed to parse config JSON: %w", err)
		}
	} else {
		config = map[string]interface{}{}
		if req.DefaultConfig {
			configFile, err := os.ReadFile("default-project-config.json")
			if err != nil {
				return nil, fmt.Errorf("failed to read default config file: %w", err)
			}
			if err := json.Unmarshal(configFile, &config); err != nil {
				return nil, fmt.Errorf("failed to parse default config JSON: %w", err)
			}
		}
	}

	return s.db.CreateOrUpdateProject(ctx, req.ID, req.Title, req.ShortName, config)
}

// DeleteProject deletes a project by ID
func (s *ProjectService) DeleteProject(ctx context.Context, id int) error {
	return s.db.Project.DeleteOneID(id).Exec(ctx)
}

// Close closes the database connection
func (s *ProjectService) Close() error {
	return s.db.Close()
}
