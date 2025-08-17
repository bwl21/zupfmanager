package models

import "time"

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string            `json:"error" example:"validation failed"`
	Message string            `json:"message,omitempty" example:"Invalid input provided"`
	Details map[string]string `json:"details,omitempty"`
} // @name ErrorResponse

// ImportFileRequest represents a file import request
type ImportFileRequest struct {
	FilePath string `json:"file_path" binding:"required" example:"/path/to/song.abc"`
} // @name ImportFileRequest

// ImportDirectoryRequest represents a directory import request
type ImportDirectoryRequest struct {
	DirectoryPath string `json:"directory_path" binding:"required" example:"/path/to/songs/"`
} // @name ImportDirectoryRequest

// ImportResult represents the result of an import operation
type ImportResult struct {
	Filename string   `json:"filename" example:"song.abc"`
	Title    string   `json:"title" example:"Amazing Grace"`
	Action   string   `json:"action" example:"created"`
	Changes  []string `json:"changes,omitempty" example:"title,genre"`
	Error    string   `json:"error,omitempty" example:"file not found"`
} // @name ImportResult

// ImportResponse represents the response from import operations
type ImportResponse struct {
	Success bool           `json:"success" example:"true"`
	Results []ImportResult `json:"results"`
	Summary ImportSummary  `json:"summary"`
} // @name ImportResponse

// ImportSummary provides statistics about the import operation
type ImportSummary struct {
	Total     int `json:"total" example:"10"`
	Created   int `json:"created" example:"7"`
	Updated   int `json:"updated" example:"2"`
	Unchanged int `json:"unchanged" example:"1"`
	Errors    int `json:"errors" example:"0"`
} // @name ImportSummary

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2025-08-16T13:00:00Z"`
	Version   string    `json:"version" example:"1.0.0"`
} // @name HealthResponse

// CreateProjectRequest represents a create project request
type CreateProjectRequest struct {
	Title         string `json:"title" binding:"required" example:"My Music Project"`
	ShortName     string `json:"short_name" binding:"required" example:"my-project"`
	ConfigFile    string `json:"config_file,omitempty" example:"/path/to/config.json"`
	DefaultConfig bool   `json:"default_config,omitempty" example:"true"`
} // @name CreateProjectRequest

// UpdateProjectRequest represents an update project request
type UpdateProjectRequest struct {
	Title         string `json:"title" binding:"required" example:"Updated Project"`
	ShortName     string `json:"short_name" binding:"required" example:"updated-project"`
	ConfigFile    string `json:"config_file,omitempty" example:"/path/to/config.json"`
	DefaultConfig bool   `json:"default_config,omitempty" example:"false"`
} // @name UpdateProjectRequest

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID        int                    `json:"id" example:"1"`
	Title     string                 `json:"title" example:"My Music Project"`
	ShortName string                 `json:"short_name" example:"my-project"`
	Config    map[string]interface{} `json:"config,omitempty"`
} // @name ProjectResponse

// ProjectListResponse represents a list of projects response
type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Count    int               `json:"count" example:"5"`
} // @name ProjectListResponse

// SongResponse represents a song response
type SongResponse struct {
	ID        int    `json:"id" example:"1"`
	Title     string `json:"title" example:"Amazing Grace"`
	Filename  string `json:"filename" example:"amazing_grace.abc"`
	Genre     string `json:"genre,omitempty" example:"Hymn"`
	Copyright string `json:"copyright,omitempty" example:"Public Domain"`
	Tocinfo   string `json:"tocinfo,omitempty" example:"John Newton"`
} // @name SongResponse

// SongListResponse represents a list of songs response
type SongListResponse struct {
	Songs []SongResponse `json:"songs"`
	Count int            `json:"count" example:"10"`
} // @name SongListResponse

// AddSongToProjectRequest represents a request to add a song to a project
type AddSongToProjectRequest struct {
	Difficulty *string `json:"difficulty,omitempty" example:"medium" enums:"easy,medium,hard,expert"`
	Priority   *int    `json:"priority,omitempty" example:"1" minimum:"1" maximum:"4"`
	Comment    *string `json:"comment,omitempty" example:"Great song for beginners"`
} // @name AddSongToProjectRequest

// UpdateProjectSongRequest represents a request to update a project-song relationship
type UpdateProjectSongRequest struct {
	Difficulty *string `json:"difficulty,omitempty" example:"hard" enums:"easy,medium,hard,expert"`
	Priority   *int    `json:"priority,omitempty" example:"2" minimum:"1" maximum:"4"`
	Comment    *string `json:"comment,omitempty" example:"Updated comment"`
} // @name UpdateProjectSongRequest

// ProjectSongResponse represents a project-song relationship
type ProjectSongResponse struct {
	ID         int              `json:"id" example:"1"`
	ProjectID  int              `json:"project_id" example:"1"`
	SongID     int              `json:"song_id" example:"1"`
	Difficulty string           `json:"difficulty" example:"medium" enums:"easy,medium,hard,expert"`
	Priority   int              `json:"priority" example:"1" minimum:"1" maximum:"4"`
	Comment    *string          `json:"comment,omitempty" example:"Great song"`
	Song       *SongResponse    `json:"song,omitempty"`
	Project    *ProjectResponse `json:"project,omitempty"`
} // @name ProjectSongResponse

// ProjectSongsResponse represents a list of project-song relationships
type ProjectSongsResponse struct {
	ProjectSongs []ProjectSongResponse `json:"project_songs"`
	Total        int                   `json:"total" example:"5"`
} // @name ProjectSongsResponse

// BuildProjectRequest represents a request to build a project
type BuildProjectRequest struct {
	OutputDir         *string `json:"output_dir,omitempty" example:"my-project"`
	AbcFileDir        *string `json:"abc_file_dir,omitempty" example:"/path/to/abc/files"`
	PriorityThreshold *int    `json:"priority_threshold,omitempty" example:"2" minimum:"1" maximum:"4"`
	SampleID          *string `json:"sample_id,omitempty" example:"sample123"`
} // @name BuildProjectRequest

// BuildStatusResponse represents the status of a build operation
type BuildStatusResponse struct {
	Status      string `json:"status" example:"running" enums:"pending,running,completed,failed"`
	Progress    int    `json:"progress" example:"75" minimum:"0" maximum:"100"`
	Message     string `json:"message,omitempty" example:"Building songs..."`
	StartedAt   string `json:"started_at,omitempty" example:"2025-08-17T18:00:00Z"`
	CompletedAt string `json:"completed_at,omitempty" example:"2025-08-17T18:05:00Z"`
	Error       string `json:"error,omitempty" example:"Build failed: file not found"`
} // @name BuildStatusResponse

// BuildResultResponse represents the result of a build operation
type BuildResultResponse struct {
	BuildID        string   `json:"build_id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ProjectID      int      `json:"project_id" example:"1"`
	Status         string   `json:"status" example:"completed" enums:"pending,running,completed,failed"`
	OutputDir      string   `json:"output_dir" example:"my-project"`
	GeneratedFiles []string `json:"generated_files,omitempty" example:"my-project/druckdateien,my-project/pdf"`
	StartedAt      string   `json:"started_at" example:"2025-08-17T18:00:00Z"`
	CompletedAt    string   `json:"completed_at,omitempty" example:"2025-08-17T18:05:00Z"`
	Error          string   `json:"error,omitempty" example:"Build failed: file not found"`
} // @name BuildResultResponse

// BuildListResponse represents a list of build results
type BuildListResponse struct {
	Builds []BuildResultResponse `json:"builds"`
	Total  int                   `json:"total" example:"3"`
} // @name BuildListResponse
