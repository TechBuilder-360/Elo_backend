package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Social holds the schema definition for the Social entity.
type Social struct {
	ent.Schema
}

// Fields of the Social.
func (Social) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.String("name").NotEmpty(),
		field.String("url").NotEmpty(),
	}
}

// Edges of the Social.
func (Social) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("socials", Business.Type).
			Ref("social").
			Unique().
			Required(),
	}
}
