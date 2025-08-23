package core

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/project"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
	"github.com/google/uuid"
)

// Common errors
var (
	ErrProjectNotFound        = errors.New("project not found")
	ErrSongNotFound          = errors.New("song not found")
	ErrSongAlreadyInProject  = errors.New("song already in project")
	ErrProjectSongNotFound   = errors.New("project-song relationship not found")
	ErrBuildNotFound         = errors.New("build not found")
	ErrBuildInProgress       = errors.New("build already in progress")
)

// In-memory build tracking (in production, this would be in database or cache)
var (
	buildStatuses = make(map[string]*BuildStatus)
	buildResults  = make(map[string]*BuildResult)
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
	var config map[string]interface{}
	var err error
	
	// Priority: direct config > config file > default config
	if len(req.Config) > 0 {
		config = req.Config
	} else {
		config, err = s.loadConfig(req.ConfigFile, req.DefaultConfig)
		if err != nil {
			return nil, err
		}
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
	var config map[string]interface{}
	var err error
	
	// Priority: direct config > config file > default config
	if len(req.Config) > 0 {
		config = req.Config
	} else {
		config, err = s.loadConfig(req.ConfigFile, req.DefaultConfig)
		if err != nil {
			return nil, err
		}
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

// ListProjectSongs lists all songs in a project with their project associations
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

	// Convert to core models
	projectSongs := ProjectSongsFromEnt(entProjectSongs)

	// Get all song IDs from this project
	songIDs := make([]int, len(projectSongs))
	for i, ps := range projectSongs {
		if ps.Song != nil {
			songIDs[i] = ps.Song.ID
		}
	}

	// Load all project associations for these songs in a single query
	if len(songIDs) > 0 {
		allProjectSongs, err := s.db.ProjectSong.Query().
			Where(projectsong.SongIDIn(songIDs...)).
			WithProject().
			All(ctx)
		
		if err != nil {
			return nil, err
		}

		// Create a map of song ID to projects
		songProjectMap := make(map[int][]*Project)
		for _, ps := range allProjectSongs {
			if project, err := ps.Edges.ProjectOrErr(); err == nil {
				songID := ps.SongID
				if songProjectMap[songID] == nil {
					songProjectMap[songID] = make([]*Project, 0)
				}
				
				// Convert ent project to core project
				coreProject := &Project{
					ID:        project.ID,
					Title:     project.Title,
					ShortName: project.ShortName,
				}
				
				// Check if project already exists in the slice
				exists := false
				for _, existingProject := range songProjectMap[songID] {
					if existingProject.ID == coreProject.ID {
						exists = true
						break
					}
				}
				
				if !exists {
					songProjectMap[songID] = append(songProjectMap[songID], coreProject)
				}
			}
		}

		// Assign project associations to songs
		for _, ps := range projectSongs {
			if ps.Song != nil {
				ps.Song.Projects = songProjectMap[ps.Song.ID]
			}
		}
	}

	return projectSongs, nil
}

// BuildProject starts a project build operation
func (s *projectService) BuildProject(ctx context.Context, req BuildProjectRequest) (*BuildResult, error) {
	// Validate input
	if err := ValidateBuildProjectRequest(req); err != nil {
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

	// Generate unique build ID
	buildID := uuid.New().String()
	
	// Set defaults
	if req.PriorityThreshold == 0 {
		req.PriorityThreshold = 4 // Include all priorities by default
	}
	if req.OutputDir == "" {
		// Get project to determine output directory
		entProject, err := s.db.Project.Get(ctx, req.ProjectID)
		if err != nil {
			return nil, err
		}
		req.OutputDir = entProject.ShortName
	}

	// Initialize build status
	now := time.Now().Format(time.RFC3339)
	buildStatus := &BuildStatus{
		Status:    "pending",
		Progress:  0,
		Message:   "Build queued",
		StartedAt: now,
	}
	buildStatuses[buildID] = buildStatus

	// Initialize build result
	buildResult := &BuildResult{
		BuildID:   buildID,
		ProjectID: req.ProjectID,
		Status:    "pending",
		OutputDir: req.OutputDir,
		StartedAt: now,
	}
	buildResults[buildID] = buildResult

	// Start build in background (simplified version using CLI command)
	go s.executeBuild(buildID, req)

	return buildResult, nil
}

// GetBuildStatus returns the current status of a build
func (s *projectService) GetBuildStatus(ctx context.Context, buildID string) (*BuildStatus, error) {
	status, exists := buildStatuses[buildID]
	if !exists {
		return nil, ErrBuildNotFound
	}
	return status, nil
}

// ListBuilds returns all builds for a project
func (s *projectService) ListBuilds(ctx context.Context, projectID int) ([]*BuildResult, error) {
	var results []*BuildResult
	for _, result := range buildResults {
		if result.ProjectID == projectID {
			results = append(results, result)
		}
	}
	return results, nil
}

// ClearBuildHistory removes all build history for a project
func (s *projectService) ClearBuildHistory(ctx context.Context, projectID int) error {
	// Remove all builds for this project from buildResults
	for buildID, result := range buildResults {
		if result.ProjectID == projectID {
			delete(buildResults, buildID)
			delete(buildStatuses, buildID)
		}
	}
	return nil
}

// executeBuild runs the actual build process in background
func (s *projectService) executeBuild(buildID string, req BuildProjectRequest) {
	status := buildStatuses[buildID]
	result := buildResults[buildID]

	// Update status to running
	status.Status = "running"
	status.Progress = 10
	status.Message = "Starting build process"
	result.Status = "running"

	// Create progress callback
	progressCallback := func(progress int, message string) {
		status.Progress = progress
		status.Message = message
		slog.Info("Build progress", "buildID", buildID, "progress", progress, "message", message)
	}

	// Execute the build directly using core logic
	ctx := context.Background()
	err := s.ExecuteProjectBuildWithProgress(ctx, req, progressCallback)

	now := time.Now().Format(time.RFC3339)
	result.CompletedAt = now
	status.CompletedAt = now

	if err != nil {
		status.Status = "failed"
		status.Progress = 100
		status.Message = "Build failed"
		result.Status = "failed"
		result.Error = fmt.Sprintf("Build failed: %v", err)
		slog.Error("Build failed", "buildID", buildID, "error", err)
	} else {
		status.Status = "completed"
		status.Progress = 100
		status.Message = "Build completed successfully"
		result.Status = "completed"
		
		// TODO: Scan output directory to get generated files
		result.GeneratedFiles = []string{} // Placeholder
		
		slog.Info("Build completed successfully", "buildID", buildID)
		result.Status = "completed"
		
		// TODO: Parse output to extract generated files
		result.GeneratedFiles = []string{
			fmt.Sprintf("%s/druckdateien", req.OutputDir),
			fmt.Sprintf("%s/pdf", req.OutputDir),
			fmt.Sprintf("%s/abc", req.OutputDir),
		}
	}
}
