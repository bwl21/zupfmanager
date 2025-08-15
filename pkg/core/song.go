package core

import (
	"context"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// SongService handles song-related operations
type SongService struct {
	db *database.Client
}

// NewSongService creates a new song service
func NewSongService() (*SongService, error) {
	db, err := database.New()
	if err != nil {
		return nil, err
	}
	return &SongService{db: db}, nil
}

// ListSongs returns all songs
func (s *SongService) ListSongs(ctx context.Context) ([]*ent.Song, error) {
	return s.db.Song.Query().All(ctx)
}

// GetSong returns a song by ID
func (s *SongService) GetSong(ctx context.Context, id int) (*ent.Song, error) {
	return s.db.Song.Get(ctx, id)
}

// SearchSongs searches for songs by title
func (s *SongService) SearchSongs(ctx context.Context, query string) ([]*ent.Song, error) {
	return s.db.Song.Query().
		Where(song.TitleContains(query)).
		All(ctx)
}

// SearchSongsAdvanced searches for songs with advanced options
func (s *SongService) SearchSongsAdvanced(ctx context.Context, query string, searchTitle, searchFilename, searchGenre bool) ([]*ent.Song, error) {
	// If no specific fields are selected, search all fields
	if !searchTitle && !searchFilename && !searchGenre {
		searchTitle = true
		searchFilename = true
		searchGenre = true
	}

	// Build predicates for search
	var predicates []predicate.Song
	if searchTitle {
		predicates = append(predicates, song.TitleContainsFold(query))
	}
	if searchFilename {
		predicates = append(predicates, song.FilenameContainsFold(query))
	}
	if searchGenre && query != "" {
		predicates = append(predicates, song.GenreContainsFold(query))
	}

	// Query songs with the search term
	query_builder := s.db.Song.Query()
	if len(predicates) > 0 {
		query_builder = query_builder.Where(song.Or(predicates...))
	}

	return query_builder.All(ctx)
}

// Close closes the database connection
func (s *SongService) Close() error {
	return s.db.Close()
}
