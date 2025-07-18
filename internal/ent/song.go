// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/bwl21/zupfmanager/internal/ent/song"
)

// Song is the model entity for the Song schema.
type Song struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Title holds the value of the "title" field.
	Title string `json:"title,omitempty"`
	// Filename holds the value of the "filename" field.
	Filename string `json:"filename,omitempty"`
	// Genre holds the value of the "genre" field.
	Genre string `json:"genre,omitempty"`
	// Copyright holds the value of the "copyright" field.
	Copyright string `json:"copyright,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the SongQuery when eager-loading is set.
	Edges        SongEdges `json:"edges"`
	selectValues sql.SelectValues
}

// SongEdges holds the relations/edges for other nodes in the graph.
type SongEdges struct {
	// ProjectSongs holds the value of the project_songs edge.
	ProjectSongs []*ProjectSong `json:"project_songs,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ProjectSongsOrErr returns the ProjectSongs value or an error if the edge
// was not loaded in eager-loading.
func (e SongEdges) ProjectSongsOrErr() ([]*ProjectSong, error) {
	if e.loadedTypes[0] {
		return e.ProjectSongs, nil
	}
	return nil, &NotLoadedError{edge: "project_songs"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Song) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case song.FieldID:
			values[i] = new(sql.NullInt64)
		case song.FieldTitle, song.FieldFilename, song.FieldGenre, song.FieldCopyright:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Song fields.
func (s *Song) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case song.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			s.ID = int(value.Int64)
		case song.FieldTitle:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field title", values[i])
			} else if value.Valid {
				s.Title = value.String
			}
		case song.FieldFilename:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field filename", values[i])
			} else if value.Valid {
				s.Filename = value.String
			}
		case song.FieldGenre:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field genre", values[i])
			} else if value.Valid {
				s.Genre = value.String
			}
		case song.FieldCopyright:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field copyright", values[i])
			} else if value.Valid {
				s.Copyright = value.String
			}
		default:
			s.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Song.
// This includes values selected through modifiers, order, etc.
func (s *Song) Value(name string) (ent.Value, error) {
	return s.selectValues.Get(name)
}

// QueryProjectSongs queries the "project_songs" edge of the Song entity.
func (s *Song) QueryProjectSongs() *ProjectSongQuery {
	return NewSongClient(s.config).QueryProjectSongs(s)
}

// Update returns a builder for updating this Song.
// Note that you need to call Song.Unwrap() before calling this method if this Song
// was returned from a transaction, and the transaction was committed or rolled back.
func (s *Song) Update() *SongUpdateOne {
	return NewSongClient(s.config).UpdateOne(s)
}

// Unwrap unwraps the Song entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (s *Song) Unwrap() *Song {
	_tx, ok := s.config.driver.(*txDriver)
	if !ok {
		panic("ent: Song is not a transactional entity")
	}
	s.config.driver = _tx.drv
	return s
}

// String implements the fmt.Stringer.
func (s *Song) String() string {
	var builder strings.Builder
	builder.WriteString("Song(")
	builder.WriteString(fmt.Sprintf("id=%v, ", s.ID))
	builder.WriteString("title=")
	builder.WriteString(s.Title)
	builder.WriteString(", ")
	builder.WriteString("filename=")
	builder.WriteString(s.Filename)
	builder.WriteString(", ")
	builder.WriteString("genre=")
	builder.WriteString(s.Genre)
	builder.WriteString(", ")
	builder.WriteString("copyright=")
	builder.WriteString(s.Copyright)
	builder.WriteByte(')')
	return builder.String()
}

// Songs is a parsable slice of Song.
type Songs []*Song
