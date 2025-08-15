package core

import (
	"context"
)

// Project represents a project domain entity
type Project struct {
	ID        int                    `json:"id"`
	Title     string                 `json:"title"`
	ShortName string                 `json:"short_name"`
	Config    map[string]interface{} `json:"config"`
}

// Song represents a song domain entity
type Song struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Filename  string `json:"filename"`
	Genre     string `json:"genre"`
	Copyright string `json:"copyright"`
	Tocinfo   string `json:"tocinfo"`
}

// CreateProjectRequest represents the data needed to create a project
type CreateProjectRequest struct {
	Title         string `json:"title" validate:"required,min=1"`
	ShortName     string `json:"short_name" validate:"required,min=1,max=50,alphanum"`
	ConfigFile    string `json:"config_file,omitempty"`
	DefaultConfig bool   `json:"default_config,omitempty"`
}

// UpdateProjectRequest represents the data needed to update a project
type UpdateProjectRequest struct {
	ID            int    `json:"id" validate:"required,min=1"`
	Title         string `json:"title" validate:"required,min=1"`
	ShortName     string `json:"short_name" validate:"required,min=1,max=50,alphanum"`
	ConfigFile    string `json:"config_file,omitempty"`
	DefaultConfig bool   `json:"default_config,omitempty"`
}

// ImportResult represents the result of an import operation
type ImportResult struct {
	Filename string   `json:"filename"`
	Title    string   `json:"title"`
	Action   string   `json:"action"` // "created", "updated", or "unchanged"
	Changes  []string `json:"changes,omitempty"`
	Error    error    `json:"error,omitempty"`
}

// SearchOptions represents search configuration
type SearchOptions struct {
	SearchTitle    bool `json:"search_title"`
	SearchFilename bool `json:"search_filename"`
	SearchGenre    bool `json:"search_genre"`
}

// ProjectService interface defines project operations
type ProjectService interface {
	Create(ctx context.Context, req CreateProjectRequest) (*Project, error)
	Update(ctx context.Context, req UpdateProjectRequest) (*Project, error)
	List(ctx context.Context) ([]*Project, error)
	Get(ctx context.Context, id int) (*Project, error)
	Delete(ctx context.Context, id int) error
}

// SongService interface defines song operations
type SongService interface {
	List(ctx context.Context) ([]*Song, error)
	Get(ctx context.Context, id int) (*Song, error)
	Search(ctx context.Context, query string) ([]*Song, error)
	SearchAdvanced(ctx context.Context, query string, options SearchOptions) ([]*Song, error)
}

// ImportService interface defines import operations
type ImportService interface {
	ImportDirectory(ctx context.Context, directory string) ([]ImportResult, error)
	ImportFile(ctx context.Context, file string) ImportResult
}

// ConfigService interface defines configuration operations
type ConfigService interface {
	LoadFromFile(path string) (map[string]interface{}, error)
	LoadDefault() (map[string]interface{}, error)
}

// FileSystemService interface defines filesystem operations
type FileSystemService interface {
	CreateProjectDirectories(shortName string) error
}
