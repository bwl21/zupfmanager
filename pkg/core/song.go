package core

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/projectsong"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// songService implements SongService interface
type songService struct {
	db *database.Client
}

// loadProjectAssociations loads project associations for the given songs
func (s *songService) loadProjectAssociations(ctx context.Context, songs []*Song) error {
	if len(songs) == 0 {
		return nil
	}
	
	songIDs := make([]int, len(songs))
	for i, song := range songs {
		songIDs[i] = song.ID
	}
	
	// Get all project-song relationships for these songs
	allProjectSongs, err := s.db.ProjectSong.Query().
		Where(projectsong.SongIDIn(songIDs...)).
		WithProject().
		All(ctx)
	
	if err != nil {
		return err
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
	for _, song := range songs {
		song.Projects = songProjectMap[song.ID]
	}
	
	return nil
}

// NewSongServiceWithDeps creates a new song service with dependencies
func NewSongServiceWithDeps(db *database.Client) SongService {
	return &songService{
		db: db,
	}
}

// List returns all songs with their project associations
func (s *songService) List(ctx context.Context) ([]*Song, error) {
	entSongs, err := s.db.Song.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	
	songs := SongsFromEnt(entSongs)
	
	// Load project associations
	if err := s.loadProjectAssociations(ctx, songs); err != nil {
		return nil, err
	}
	
	return songs, nil
}

// Get returns a song by ID
func (s *songService) Get(ctx context.Context, id int) (*Song, error) {
	entSong, err := s.db.Song.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return SongFromEnt(entSong), nil
}

// Search searches for songs by title and filename
func (s *songService) Search(ctx context.Context, query string) ([]*Song, error) {
	entSongs, err := s.db.Song.Query().
		Where(song.Or(
			song.TitleContains(query),
			song.FilenameContains(query),
		)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	
	songs := SongsFromEnt(entSongs)
	
	// Load project associations
	if err := s.loadProjectAssociations(ctx, songs); err != nil {
		return nil, err
	}
	
	return songs, nil
}

// SearchAdvanced performs advanced search with options
func (s *songService) SearchAdvanced(ctx context.Context, query string, options SearchOptions) ([]*Song, error) {
	var predicates []predicate.Song
	
	if options.SearchTitle {
		predicates = append(predicates, song.TitleContains(query))
	}
	if options.SearchFilename {
		predicates = append(predicates, song.FilenameContains(query))
	}
	if options.SearchGenre {
		predicates = append(predicates, song.GenreContains(query))
	}
	
	if len(predicates) == 0 {
		// Default to title search if no options specified
		predicates = append(predicates, song.TitleContains(query))
	}
	
	entSongs, err := s.db.Song.Query().
		Where(song.Or(predicates...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	
	songs := SongsFromEnt(entSongs)
	
	// Load project associations
	if err := s.loadProjectAssociations(ctx, songs); err != nil {
		return nil, err
	}
	
	return songs, nil
}

// Delete removes a song from the database
func (s *songService) Delete(ctx context.Context, id int) error {
	// Check if song exists
	entSong, err := s.db.Song.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get song: %w", err)
	}
	
	// Check if song is used in any projects
	projectSongs, err := s.db.ProjectSong.Query().
		Where(projectsong.SongID(id)).
		All(ctx)
	if err != nil {
		return fmt.Errorf("failed to check project associations: %w", err)
	}
	
	if len(projectSongs) > 0 {
		return fmt.Errorf("cannot delete song '%s': it is used in %d project(s). Remove it from all projects first", entSong.Title, len(projectSongs))
	}
	
	// Delete the song
	err = s.db.Song.DeleteOneID(id).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}
	
	return nil
}

// GeneratePreview searches for existing PDFs for a song in the ABC directory
func (s *songService) GeneratePreview(ctx context.Context, req GeneratePreviewRequest) (*GeneratePreviewResponse, error) {
	// Get the song
	entSong, err := s.db.Song.Get(ctx, req.SongID)
	if err != nil {
		return nil, fmt.Errorf("failed to get song: %w", err)
	}
	
	// Check if ABC file directory exists
	if _, err := os.Stat(req.AbcFileDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("ABC file directory does not exist: %s", req.AbcFileDir)
	}
	
	// Find existing PDF files for this song
	baseFilename := strings.TrimSuffix(entSong.Filename, ".abc")
	pattern := filepath.Join(req.AbcFileDir, baseFilename+"*.pdf")
	pdfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to search for PDF files: %w", err)
	}
	
	// Extract just the filenames
	var filenames []string
	for _, pdfFile := range pdfFiles {
		filenames = append(filenames, filepath.Base(pdfFile))
	}
	
	return &GeneratePreviewResponse{
		PDFFiles:   filenames,
		PreviewDir: req.AbcFileDir, // Return the ABC directory as the "preview dir"
	}, nil
}

// ListPreviewPDFs lists available preview PDFs for a song (requires abc_file_dir to be set)
func (s *songService) ListPreviewPDFs(ctx context.Context, songID int) ([]*PreviewPDF, error) {
	// This method now requires the ABC file directory to be provided via a different approach
	// Since we don't store the directory path, we return empty list
	// The frontend should use GeneratePreview (which is now "FindPreview") to get PDFs
	return []*PreviewPDF{}, nil
}

// GetPreviewPDF returns the file path for a preview PDF (requires abc_file_dir context)
func (s *songService) GetPreviewPDF(ctx context.Context, songID int, filename string) (string, error) {
	// This method now needs to be called with the ABC directory context
	// Since we don't store the directory path, this will need to be refactored
	// to accept the directory as a parameter
	return "", fmt.Errorf("GetPreviewPDF requires ABC directory context - use GetPreviewPDFFromDir instead")
}

// GetPreviewPDFFromDir returns the file path for a preview PDF from a specific directory
func (s *songService) GetPreviewPDFFromDir(ctx context.Context, songID int, filename string, abcFileDir string) (string, error) {
	// Validate filename to prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("invalid filename")
	}
	
	// Get the song to validate the filename belongs to this song
	entSong, err := s.db.Song.Get(ctx, songID)
	if err != nil {
		return "", fmt.Errorf("failed to get song: %w", err)
	}
	
	baseFilename := strings.TrimSuffix(entSong.Filename, ".abc")
	if !strings.HasPrefix(filename, baseFilename) {
		return "", fmt.Errorf("filename does not belong to this song")
	}
	
	filePath := filepath.Join(abcFileDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF file not found")
	}
	
	return filePath, nil
}

// CleanupPreviewPDFs is not applicable when PDFs are stored in the ABC directory
func (s *songService) CleanupPreviewPDFs(ctx context.Context, songID int) error {
	// Since PDFs are now stored in the ABC directory alongside the ABC files,
	// we don't clean them up automatically as they are part of the user's workflow
	return fmt.Errorf("cleanup not supported - PDFs are stored in ABC directory")
}

