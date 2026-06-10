package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Social holds the schema definition for the Social entity.
type Social struct {
	ent.Schema
}

func (Social) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Social.
func (Social) Fields() []ent.Field {
	return []ent.Field{
		field.String("business_id"),
		field.String("name").NotEmpty(),
		field.String("url").NotEmpty(),
	}
}

// Edges of the Social.
func (Social) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("social", Business.Type).
			Ref("socials").
			Field("business_id").
			Required().
			Unique(),
	}
}
