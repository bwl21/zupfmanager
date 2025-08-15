package core

import (
	"context"

	"github.com/bwl21/zupfmanager/internal/database"
)

// projectService implements ProjectService interface
type projectService struct {
	db         *database.Client
	config     ConfigService
	fileSystem FileSystemService
}

// NewProjectServiceWithDeps creates a new project service with dependencies
func NewProjectServiceWithDeps(db *database.Client, config ConfigService, fileSystem FileSystemService) ProjectService {
	return &projectService{
		db:         db,
		config:     config,
		fileSystem: fileSystem,
	}
}

// Create creates a new project with the given parameters
func (s *projectService) Create(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	// Validate input
	if err := ValidateCreateProjectRequest(req); err != nil {
		return nil, err
	}

	// Load configuration
	config, err := s.loadConfig(req.ConfigFile, req.DefaultConfig)
	if err != nil {
		return nil, err
	}

	// Create project in database
	entProject, err := s.db.CreateOrUpdateProject(ctx, 0, req.Title, req.ShortName, config)
	if err != nil {
		return nil, err
	}

	// Create directory structure
	if err := s.fileSystem.CreateProjectDirectories(req.ShortName); err != nil {
		return nil, err
	}

	return ProjectFromEnt(entProject), nil
}

// Update updates an existing project
func (s *projectService) Update(ctx context.Context, req UpdateProjectRequest) (*Project, error) {
	// Validate input
	if err := ValidateUpdateProjectRequest(req); err != nil {
		return nil, err
	}

	// Load configuration
	config, err := s.loadConfig(req.ConfigFile, req.DefaultConfig)
	if err != nil {
		return nil, err
	}

	// Update project in database
	entProject, err := s.db.CreateOrUpdateProject(ctx, req.ID, req.Title, req.ShortName, config)
	if err != nil {
		return nil, err
	}

	return ProjectFromEnt(entProject), nil
}

// List returns all projects
func (s *projectService) List(ctx context.Context) ([]*Project, error) {
	entProjects, err := s.db.Project.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return ProjectsFromEnt(entProjects), nil
}

// Get returns a project by ID
func (s *projectService) Get(ctx context.Context, id int) (*Project, error) {
	entProject, err := s.db.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	return ProjectFromEnt(entProject), nil
}

// Delete deletes a project by ID
func (s *projectService) Delete(ctx context.Context, id int) error {
	return s.db.Project.DeleteOneID(id).Exec(ctx)
}

// loadConfig loads configuration based on request parameters
func (s *projectService) loadConfig(configFile string, useDefault bool) (map[string]interface{}, error) {
	if configFile != "" {
		return s.config.LoadFromFile(configFile)
	}
	
	if useDefault {
		return s.config.LoadDefault()
	}
	
	return map[string]interface{}{}, nil
}
