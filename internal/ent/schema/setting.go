package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Setting holds the schema definition for the Setting entity.
type Setting struct {
	ent.Schema
}

// Fields of the Setting.
func (Setting) Fields() []ent.Field {
	return []ent.Field{
		field.String("key").
			Unique().
			NotEmpty().
			Comment("Setting key (e.g., 'last_import_path')"),
		field.String("value").
			Optional().
			Comment("Setting value"),
	}
}

// Edges of the Setting.
func (Setting) Edges() []ent.Edge {
	return nil
}
