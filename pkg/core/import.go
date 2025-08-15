package core

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// importService implements ImportService interface
type importService struct {
	db *database.Client
}

// NewImportServiceWithDeps creates a new import service with dependencies
func NewImportServiceWithDeps(db *database.Client) ImportService {
	return &importService{
		db: db,
	}
}

// ImportDirectory imports all ABC files from a directory
func (s *importService) ImportDirectory(ctx context.Context, directory string) ([]ImportResult, error) {
	files, err := filepath.Glob(filepath.Join(directory, "*.abc"))
	if err != nil {
		return nil, err
	}

	results := make([]ImportResult, 0, len(files))
	for _, file := range files {
		result := s.ImportFile(ctx, file)
		results = append(results, result)
	}

	return results, nil
}

// ImportFile imports a single ABC file
func (s *importService) ImportFile(ctx context.Context, file string) ImportResult {
	filename := filepath.Base(file)
	result := ImportResult{Filename: filename}

	content, err := os.ReadFile(file)
	if err != nil {
		result.Error = err
		return result
	}

	metadata := s.parseABCMetadata(content)
	if metadata.Title == "" {
		result.Error = fmt.Errorf("no title found in file")
		return result
	}

	result.Title = metadata.Title

	// Check if song already exists
	existingSong, err := s.db.Song.Query().Where(song.Filename(filename)).First(ctx)
	if err != nil && !ent.IsNotFound(err) {
		result.Error = fmt.Errorf("failed to query song: %w", err)
		return result
	}

	if existingSong == nil {
		// Create new song
		_, err = s.db.Song.Create().
			SetTitle(metadata.Title).
			SetFilename(filename).
			SetGenre(metadata.Genre).
			SetCopyright(metadata.Copyright).
			SetTocinfo(metadata.Tocinfo).
			Save(ctx)
		if err != nil {
			result.Error = fmt.Errorf("failed to create song: %w", err)
			return result
		}
		result.Action = "created"
	} else {
		// Update existing song
		changes := s.detectChanges(existingSong, metadata)
		if len(changes) > 0 {
			_, err = existingSong.Update().
				SetTitle(metadata.Title).
				SetGenre(metadata.Genre).
				SetCopyright(metadata.Copyright).
				SetTocinfo(metadata.Tocinfo).
				Save(ctx)
			if err != nil {
				result.Error = fmt.Errorf("failed to update song: %w", err)
				return result
			}
			result.Action = "updated"
			result.Changes = changes
		} else {
			result.Action = "unchanged"
		}
	}

	return result
}

// ABCMetadata represents metadata extracted from an ABC file
type ABCMetadata struct {
	Title     string
	Genre     string
	Copyright string
	Tocinfo   string
}

// parseABCMetadata extracts metadata from ABC file content
func (s *importService) parseABCMetadata(content []byte) ABCMetadata {
	var metadata ABCMetadata

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "T:") {
			metadata.Title = strings.TrimSpace(strings.TrimPrefix(line, "T:"))
		} else if strings.HasPrefix(line, "Z:genre") {
			metadata.Genre = strings.TrimSpace(strings.TrimPrefix(line, "Z:genre"))
		} else if strings.HasPrefix(line, "Z:copyright") {
			metadata.Copyright = strings.TrimSpace(strings.TrimPrefix(line, "Z:copyright"))
		} else if metadata.Tocinfo == "" {
			// Check for C:M: patterns
			if strings.HasPrefix(line, "C:M: ") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M: "))
			} else if strings.HasPrefix(line, "C:M:") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M:"))
			} else if strings.HasPrefix(line, "C:M+T: ") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M+T: "))
			} else if strings.HasPrefix(line, "C:M+T:") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:M+T:"))
			} else if strings.HasPrefix(line, "C:T+M: ") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:T+M: "))
			} else if strings.HasPrefix(line, "C:T+M:") {
				metadata.Tocinfo = strings.TrimSpace(strings.TrimPrefix(line, "C:T+M:"))
			}
		}
	}

	return metadata
}

// detectChanges compares existing song with new metadata and returns list of changes
func (s *importService) detectChanges(existing *ent.Song, metadata ABCMetadata) []string {
	changes := make([]string, 0)

	if existing.Title != metadata.Title {
		changes = append(changes, fmt.Sprintf("title: %s -> %s", existing.Title, metadata.Title))
	}
	if existing.Genre != metadata.Genre {
		changes = append(changes, fmt.Sprintf("genre: %s -> %s", existing.Genre, metadata.Genre))
	}
	if existing.Copyright != metadata.Copyright {
		changes = append(changes, fmt.Sprintf("copyright: %s -> %s", existing.Copyright, metadata.Copyright))
	}
	if existing.Tocinfo != metadata.Tocinfo {
		changes = append(changes, fmt.Sprintf("tocinfo: %s -> %s", existing.Tocinfo, metadata.Tocinfo))
	}

	return changes
}
