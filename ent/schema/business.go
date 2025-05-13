package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Toflex/directory_v2/pkg/utils"
)

// Business holds the schema definition for the Business entity.
type Business struct {
	ent.Schema
}

// Fields of the Business.
func (Business) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").Default("others"),
		field.String("name").Unique().NotEmpty(),
		field.String("logo_url").Nillable(),
		field.String("email").NotEmpty().
			Validate(utils.ValidateEmail),
		field.String("website").Nillable(),
		field.Bool("disabled").Default(true),
		field.Time("disabled_at").Default(time.Now),
		field.String("disable_reason").Nillable(),
		field.Bool("verified").Default(false),
		field.Time("verified_at").Nillable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Business.
func (Business) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("social", Social.Type),
		edge.To("manager", Manager.Type),
	}
}
