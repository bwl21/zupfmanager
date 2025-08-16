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
