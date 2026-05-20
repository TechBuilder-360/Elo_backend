package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Provider holds the schema definition for the Provider entity.
type Provider struct {
	ent.Schema
}

// Fields of the Provider.
func (Provider) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),
		field.String("slug").
			Unique().
			NotEmpty(),
		field.Bool("active"),
	}
}

// Edges of the Provider.
func (Provider) Edges() []ent.Edge {
	return nil
}
