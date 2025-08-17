package core

import (
	"context"
	"errors"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// Common errors
var (
	ErrProjectNotFound        = errors.New("project not found")
	ErrSongNotFound          = errors.New("song not found")
	ErrSongAlreadyInProject  = errors.New("song already in project")
	ErrProjectSongNotFound   = errors.New("project-song relationship not found")
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

// AddSongToProject adds a song to a project with optional difficulty, priority, and comment
func (s *projectService) AddSongToProject(ctx context.Context, req AddSongToProjectRequest) (*ProjectSong, error) {
	// Validate input
	if err := ValidateAddSongToProjectRequest(req); err != nil {
		return nil, err
	}

	// Check if project exists
	projectExists, err := s.db.Project.Query().Where(project.ID(req.ProjectID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !projectExists {
		return nil, ErrProjectNotFound
	}

	// Check if song exists
	songExists, err := s.db.Song.Query().Where(song.ID(req.SongID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !songExists {
		return nil, ErrSongNotFound
	}

	// Check if song is already in project
	exists, err := s.db.ProjectSong.Query().
		Where(
			projectsong.ProjectID(req.ProjectID),
			projectsong.SongID(req.SongID),
		).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSongAlreadyInProject
	}

	// Create project-song relationship
	builder := s.db.ProjectSong.Create().
		SetProjectID(req.ProjectID).
		SetSongID(req.SongID)

	if req.Difficulty != nil {
		builder = builder.SetDifficulty(projectsong.Difficulty(*req.Difficulty))
	}
	if req.Priority != nil {
		builder = builder.SetPriority(*req.Priority)
	}
	if req.Comment != nil {
		builder = builder.SetComment(*req.Comment)
	}

	entProjectSong, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return ProjectSongFromEnt(entProjectSong), nil
}

// RemoveSongFromProject removes a song from a project
func (s *projectService) RemoveSongFromProject(ctx context.Context, projectID, songID int) error {
	deleted, err := s.db.ProjectSong.Delete().
		Where(
			projectsong.ProjectID(projectID),
			projectsong.SongID(songID),
		).
		Exec(ctx)

	if err != nil {
		return err
	}

	if deleted == 0 {
		return ErrProjectSongNotFound
	}

	return nil
}

// UpdateProjectSong updates a project-song relationship
func (s *projectService) UpdateProjectSong(ctx context.Context, req UpdateProjectSongRequest) (*ProjectSong, error) {
	// Validate input
	if err := ValidateUpdateProjectSongRequest(req); err != nil {
		return nil, err
	}

	// Find the project-song relationship
	entProjectSong, err := s.db.ProjectSong.Query().
		Where(
			projectsong.ProjectID(req.ProjectID),
			projectsong.SongID(req.SongID),
		).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrProjectSongNotFound
		}
		return nil, err
	}

	// Update the relationship
	builder := entProjectSong.Update()
	if req.Difficulty != nil {
		builder = builder.SetDifficulty(projectsong.Difficulty(*req.Difficulty))
	}
	if req.Priority != nil {
		builder = builder.SetPriority(*req.Priority)
	}
	if req.Comment != nil {
		builder = builder.SetComment(*req.Comment)
	}

	updatedProjectSong, err := builder.Save(ctx)
	if err != nil {
		return nil, err
	}

	return ProjectSongFromEnt(updatedProjectSong), nil
}

// ListProjectSongs lists all songs in a project
func (s *projectService) ListProjectSongs(ctx context.Context, projectID int) ([]*ProjectSong, error) {
	// Check if project exists
	projectExists, err := s.db.Project.Query().Where(project.ID(projectID)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if !projectExists {
		return nil, ErrProjectNotFound
	}

	// Get all project-song relationships for this project
	entProjectSongs, err := s.db.ProjectSong.Query().
		Where(projectsong.ProjectID(projectID)).
		WithSong().
		WithProject().
		All(ctx)

	if err != nil {
		return nil, err
	}

	return ProjectSongsFromEnt(entProjectSongs), nil
}
