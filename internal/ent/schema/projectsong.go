package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ProjectSong holds the schema definition for the ProjectSong entity.
type ProjectSong struct {
	ent.Schema
}

// Fields of the ProjectSong.
func (ProjectSong) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable().
			StructTag(`json:"id,omitempty"`),
		field.Int("priority").
			Range(1, 4),
		field.Enum("difficulty").
			Values("easy", "medium", "hard", "expert").
			Default("medium"),
		field.String("comment").
			Optional(),
		field.Int("project_id"),
		field.Int("song_id"),
	}
}

// Edges of the ProjectSong.
func (ProjectSong) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project", Project.Type).
			Field("project_id").
			Unique().
			Required(),
		edge.To("song", Song.Type).
			Field("song_id").
			Unique().
			Required(),
	}
}

// Indexes of the ProjectSong.
func (ProjectSong) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("project_id", "song_id").
			Unique(),
	}
}
