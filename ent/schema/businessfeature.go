package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// BusinessFeature holds the schema definition for the BusinessFeature entity.
type BusinessFeature struct {
	ent.Schema
}

func (BusinessFeature) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the BusinessFeature.
func (BusinessFeature) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("identifier").
			NotEmpty().
			Unique(),
		field.Bool("require_subscription").
			Comment("If set to true, businesses will need to subscribe to have access to this feature, else it is generally available to all business").
			Default(false),
		field.Bool("active").
			Default(true),
		field.Int("min").
			Default(0),
		field.Int("max").
			Default(0),
		field.JSON("fee", &Fee{}),
	}
}

// Edges of the BusinessFeature.
func (BusinessFeature) Edges() []ent.Edge {
	return nil
}
