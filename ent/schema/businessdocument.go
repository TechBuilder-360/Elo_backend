package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BusinessDocument holds the schema definition for the BusinessDocument entity.
type BusinessDocument struct {
	ent.Schema
}

func (BusinessDocument) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the BusinessDocument.
func (BusinessDocument) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("description").NotEmpty(),
		field.String("url").NotEmpty(),
		field.Bool("verified").Default(false).Comment("this field only applies to document_type KYB"),
		field.Enum("type").Values(
			"KYB",
			"SERVICE",
			"PRODUCT",
		),
	}
}

// Edges of the BusinessDocument.
func (BusinessDocument) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("business", Business.Type).
			Ref("business_documents").
			Required().
			Unique(),
		edge.From("kyb_document", KYBDocument.Type).
			Ref("kyb_documents"),
	}
}
