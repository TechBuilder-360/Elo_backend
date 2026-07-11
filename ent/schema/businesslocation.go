package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BusinessLocation holds the schema definition for the BusinessLocation entity.
type BusinessLocation struct {
	ent.Schema
}

func (BusinessLocation) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the BusinessLocation.
func (BusinessLocation) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_head_office").Default(false),
		field.String("name").NotEmpty(),
		field.String("address").NotEmpty(),
		field.String("city").NotEmpty(),
		field.String("state").NotEmpty(),
		field.String("country").NotEmpty(),
		field.String("zip_code").NotEmpty(),
		field.Float("latitude").Optional(),
		field.Float("longitude").Optional(),
		field.Bool("active").Default(true),
	}
}

// Edges of the BusinessLocation.
func (BusinessLocation) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("business", Business.Type).
			Ref("locations").
			Unique(),
	}
}
