package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/song"
	"github.com/bwl21/zupfmanager/internal/zupfnoter"
)

// songService implements SongService interface
type songService struct {
	db *database.Client
}

// NewSongServiceWithDeps creates a new song service with dependencies
func NewSongServiceWithDeps(db *database.Client) SongService {
	return &songService{
		db: db,
	}
}

// List returns all songs
func (s *songService) List(ctx context.Context) ([]*Song, error) {
	entSongs, err := s.db.Song.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return SongsFromEnt(entSongs), nil
}

// Get returns a song by ID
func (s *songService) Get(ctx context.Context, id int) (*Song, error) {
	entSong, err := s.db.Song.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return SongFromEnt(entSong), nil
}

// Search searches for songs by title
func (s *songService) Search(ctx context.Context, query string) ([]*Song, error) {
	entSongs, err := s.db.Song.Query().
		Where(song.TitleContains(query)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return SongsFromEnt(entSongs), nil
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
	return SongsFromEnt(entSongs), nil
}

// GeneratePreview generates preview PDFs for a song
func (s *songService) GeneratePreview(ctx context.Context, req GeneratePreviewRequest) (*GeneratePreviewResponse, error) {
	// Get the song
	entSong, err := s.db.Song.Get(ctx, req.SongID)
	if err != nil {
		return nil, fmt.Errorf("failed to get song: %w", err)
	}
	
	// Create preview directory
	previewDir := filepath.Join(os.TempDir(), "zupfmanager", "previews", fmt.Sprintf("song-%d", req.SongID))
	err = os.MkdirAll(previewDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create preview directory: %w", err)
	}
	
	// Create PDF output directory
	pdfDir := filepath.Join(previewDir, "pdf")
	err = os.MkdirAll(pdfDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create PDF directory: %w", err)
	}
	
	// Create config file if provided
	var configFile string
	if req.Config != nil && len(req.Config) > 0 {
		tempConfigFile, err := os.CreateTemp("", "zupfnoter-preview-*.json")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp config file: %w", err)
		}
		defer os.Remove(tempConfigFile.Name())
		
		err = json.NewEncoder(tempConfigFile).Encode(req.Config)
		if err != nil {
			return nil, fmt.Errorf("failed to encode config: %w", err)
		}
		tempConfigFile.Close()
		configFile = tempConfigFile.Name()
	}
	
	// Run zupfnoter to generate PDFs
	abcFilePath := filepath.Join(req.AbcFileDir, entSong.Filename)
	var args []string
	if configFile != "" {
		args = []string{abcFilePath, pdfDir, configFile}
	} else {
		args = []string{abcFilePath, pdfDir}
	}
	
	_, _, err = zupfnoter.Run(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("zupfnoter failed: %w", err)
	}
	
	// Find generated PDF files
	baseFilename := strings.TrimSuffix(entSong.Filename, ".abc")
	pattern := filepath.Join(pdfDir, baseFilename+"*.pdf")
	pdfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find generated PDFs: %w", err)
	}
	
	// Extract just the filenames
	var filenames []string
	for _, pdfFile := range pdfFiles {
		filenames = append(filenames, filepath.Base(pdfFile))
	}
	
	return &GeneratePreviewResponse{
		PDFFiles:   filenames,
		PreviewDir: previewDir,
	}, nil
}

// ListPreviewPDFs lists available preview PDFs for a song
func (s *songService) ListPreviewPDFs(ctx context.Context, songID int) ([]*PreviewPDF, error) {
	previewDir := filepath.Join(os.TempDir(), "zupfmanager", "previews", fmt.Sprintf("song-%d", songID), "pdf")
	
	// Check if preview directory exists
	if _, err := os.Stat(previewDir); os.IsNotExist(err) {
		return []*PreviewPDF{}, nil
	}
	
	// Find PDF files
	pattern := filepath.Join(previewDir, "*.pdf")
	pdfFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find PDF files: %w", err)
	}
	
	var previews []*PreviewPDF
	for _, pdfFile := range pdfFiles {
		stat, err := os.Stat(pdfFile)
		if err != nil {
			continue // Skip files we can't stat
		}
		
		previews = append(previews, &PreviewPDF{
			Filename:  filepath.Base(pdfFile),
			Size:      stat.Size(),
			CreatedAt: stat.ModTime().Format(time.RFC3339),
		})
	}
	
	return previews, nil
}

// GetPreviewPDF returns the file path for a preview PDF
func (s *songService) GetPreviewPDF(ctx context.Context, songID int, filename string) (string, error) {
	// Validate filename to prevent path traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return "", fmt.Errorf("invalid filename")
	}
	
	previewDir := filepath.Join(os.TempDir(), "zupfmanager", "previews", fmt.Sprintf("song-%d", songID), "pdf")
	filePath := filepath.Join(previewDir, filename)
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("PDF file not found")
	}
	
	return filePath, nil
}

// CleanupPreviewPDFs removes all preview PDFs for a song
func (s *songService) CleanupPreviewPDFs(ctx context.Context, songID int) error {
	previewDir := filepath.Join(os.TempDir(), "zupfmanager", "previews", fmt.Sprintf("song-%d", songID))
	
	if _, err := os.Stat(previewDir); os.IsNotExist(err) {
		return nil // Nothing to clean up
	}
	
	return os.RemoveAll(previewDir)
}

