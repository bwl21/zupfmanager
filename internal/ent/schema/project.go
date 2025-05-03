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
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("project_songs", ProjectSong.Type).
			Ref("project"),
	}
}
