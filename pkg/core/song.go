package core

import (
	"context"

	"github.com/bwl21/zupfmanager/internal/database"
	"github.com/bwl21/zupfmanager/internal/ent/predicate"
	"github.com/bwl21/zupfmanager/internal/ent/song"
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

// SearchAdvanced searches for songs with advanced options
func (s *songService) SearchAdvanced(ctx context.Context, query string, options SearchOptions) ([]*Song, error) {
	// If no specific fields are selected, search all fields
	if !options.SearchTitle && !options.SearchFilename && !options.SearchGenre {
		options.SearchTitle = true
		options.SearchFilename = true
		options.SearchGenre = true
	}

	// Build predicates for search
	var predicates []predicate.Song
	if options.SearchTitle {
		predicates = append(predicates, song.TitleContainsFold(query))
	}
	if options.SearchFilename {
		predicates = append(predicates, song.FilenameContainsFold(query))
	}
	if options.SearchGenre && query != "" {
		predicates = append(predicates, song.GenreContainsFold(query))
	}

	// Query songs with the search term
	queryBuilder := s.db.Song.Query()
	if len(predicates) > 0 {
		queryBuilder = queryBuilder.Where(song.Or(predicates...))
	}

	entSongs, err := queryBuilder.All(ctx)
	if err != nil {
		return nil, err
	}
	return SongsFromEnt(entSongs), nil
}
