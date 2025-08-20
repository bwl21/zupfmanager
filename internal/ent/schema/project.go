package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Positive().
			Immutable().
			StructTag(`json:"id,omitempty"`),
		field.String("title").
			NotEmpty(),
		field.String("short_name").
			NotEmpty(),
		field.JSON("config", map[string]interface{}{}).
			Optional(),
		field.String("abc_file_dir_preference").
			Optional().
			Comment("User's preferred directory for ABC files"),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("project_songs", ProjectSong.Type),
	}
}
