package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// KYBDocument holds the schema definition for the KYBDocument entity.
type KYBDocument struct {
	ent.Schema
}

func (KYBDocument) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the KYBDocument.
func (KYBDocument) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.Bool("required").Default(true),
		field.Bool("active").Default(true),
	}
}

// Edges of the KYBDocument.
func (KYBDocument) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("kyb_documents", BusinessDocument.Type),
	}
}
