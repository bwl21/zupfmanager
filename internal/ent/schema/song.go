package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Song holds the schema definition for the Song entity.
type Song struct {
	ent.Schema
}

// Fields of the Song.
func (Song) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable().
			StructTag(`json:"id,omitempty"`),
		field.String("title").
			NotEmpty(),
		field.String("filename").
			NotEmpty().
			Unique(),
		field.String("genre").
			Optional(),
		field.String("copyright").
			Optional(),
	}
}

func (Song) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("filename").Unique(),
		index.Fields("title"),
		index.Fields("genre"),
	}
}

// Edges of the Song.
func (Song) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project_songs", ProjectSong.Type).
			Ref("song"),
	}
}
