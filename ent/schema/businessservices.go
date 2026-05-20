package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BusinessServices holds the schema definition for the BusinessServices entity.
type BusinessServices struct {
	ent.Schema
}

// Fields of the BusinessServices.
func (BusinessServices) Fields() []ent.Field {
	return []ent.Field{
		field.String("business_id"),
		field.String("title").NotEmpty(),
		field.String("description").NotEmpty(),
	}
}

// Edges of the BusinessServices.
func (BusinessServices) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("business_service", Business.Type).
			Ref("services").
			Field("business_id").
			Required().
			Unique(),
	}
}
