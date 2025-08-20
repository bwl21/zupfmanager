package core

import (
	"context"
)

// Project represents a project domain entity
type Project struct {
	ID                   int                    `json:"id"`
	Title                string                 `json:"title"`
	ShortName            string                 `json:"short_name"`
	Config               map[string]interface{} `json:"config"`
	AbcFileDirPreference string                 `json:"abc_file_dir_preference,omitempty"`
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

// ProjectSong represents a project-song relationship
type ProjectSong struct {
	ID         int     `json:"id"`
	ProjectID  int     `json:"project_id"`
	SongID     int     `json:"song_id"`
	Difficulty string  `json:"difficulty"`
	Priority   int     `json:"priority"`
	Comment    *string `json:"comment,omitempty"`
	Song       *Song   `json:"song,omitempty"`
	Project    *Project `json:"project,omitempty"`
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

// AddSongToProjectRequest represents the data needed to add a song to a project
type AddSongToProjectRequest struct {
	ProjectID  int     `json:"project_id" validate:"required,min=1"`
	SongID     int     `json:"song_id" validate:"required,min=1"`
	Difficulty *string `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard expert"`
	Priority   *int    `json:"priority,omitempty" validate:"omitempty,min=1,max=4"`
	Comment    *string `json:"comment,omitempty"`
}

// UpdateProjectSongRequest represents the data needed to update a project-song relationship
type UpdateProjectSongRequest struct {
	ProjectID  int     `json:"project_id" validate:"required,min=1"`
	SongID     int     `json:"song_id" validate:"required,min=1"`
	Difficulty *string `json:"difficulty,omitempty" validate:"omitempty,oneof=easy medium hard expert"`
	Priority   *int    `json:"priority,omitempty" validate:"omitempty,min=1,max=4"`
	Comment    *string `json:"comment,omitempty"`
}

// BuildProjectRequest represents the data needed to build a project
type BuildProjectRequest struct {
	ProjectID         int    `json:"project_id" validate:"required,min=1"`
	OutputDir         string `json:"output_dir,omitempty"`
	AbcFileDir        string `json:"abc_file_dir,omitempty"`
	PriorityThreshold int    `json:"priority_threshold,omitempty" validate:"omitempty,min=1,max=4"`
	SampleID          string `json:"sample_id,omitempty"`
}

// BuildStatus represents the status of a build operation
type BuildStatus struct {
	Status      string `json:"status"` // "pending", "running", "completed", "failed"
	Progress    int    `json:"progress"` // 0-100
	Message     string `json:"message,omitempty"`
	StartedAt   string `json:"started_at,omitempty"`
	CompletedAt string `json:"completed_at,omitempty"`
	Error       string `json:"error,omitempty"`
}

// BuildResult represents the result of a build operation
type BuildResult struct {
	BuildID     string   `json:"build_id"`
	ProjectID   int      `json:"project_id"`
	Status      string   `json:"status"`
	OutputDir   string   `json:"output_dir"`
	GeneratedFiles []string `json:"generated_files,omitempty"`
	StartedAt   string   `json:"started_at"`
	CompletedAt string   `json:"completed_at,omitempty"`
	Error       string   `json:"error,omitempty"`
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
	
	// Project-Song relationship operations
	AddSongToProject(ctx context.Context, req AddSongToProjectRequest) (*ProjectSong, error)
	RemoveSongFromProject(ctx context.Context, projectID, songID int) error
	UpdateProjectSong(ctx context.Context, req UpdateProjectSongRequest) (*ProjectSong, error)
	ListProjectSongs(ctx context.Context, projectID int) ([]*ProjectSong, error)
	
	// Project build operations
	BuildProject(ctx context.Context, req BuildProjectRequest) (*BuildResult, error)
	ExecuteProjectBuild(ctx context.Context, req BuildProjectRequest) error
	GetBuildStatus(ctx context.Context, buildID string) (*BuildStatus, error)
	ListBuilds(ctx context.Context, projectID int) ([]*BuildResult, error)
}

// PreviewPDF represents a generated preview PDF
type PreviewPDF struct {
	Filename  string `json:"filename"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"created_at"`
}

// GeneratePreviewRequest represents the request to generate preview PDFs
type GeneratePreviewRequest struct {
	SongID     int                    `json:"song_id" validate:"required,min=1"`
	AbcFileDir string                 `json:"abc_file_dir" validate:"required"`
	Config     map[string]interface{} `json:"config,omitempty"`
}

// GeneratePreviewResponse represents the response from preview generation
type GeneratePreviewResponse struct {
	PDFFiles   []string `json:"pdf_files"`
	PreviewDir string   `json:"preview_dir"`
}

// SongService interface defines song operations
type SongService interface {
	List(ctx context.Context) ([]*Song, error)
	Get(ctx context.Context, id int) (*Song, error)
	Search(ctx context.Context, query string) ([]*Song, error)
	SearchAdvanced(ctx context.Context, query string, options SearchOptions) ([]*Song, error)
	
	// Preview operations
	GeneratePreview(ctx context.Context, req GeneratePreviewRequest) (*GeneratePreviewResponse, error)
	ListPreviewPDFs(ctx context.Context, songID int) ([]*PreviewPDF, error)
	GetPreviewPDF(ctx context.Context, songID int, filename string) (string, error) // returns file path
	CleanupPreviewPDFs(ctx context.Context, songID int) error
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
