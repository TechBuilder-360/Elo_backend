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

func (Business) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Business.
func (Business) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").Default("others"),
		field.String("name").Unique().NotEmpty(),
		field.String("logo").Optional(),
		field.String("email").NotEmpty().
			Validate(utils.ValidateEmail),
		field.String("website").Optional(),
		field.Bool("disabled").Default(true),
		field.Time("disabled_at").Default(time.Now),
		field.String("disable_reason").Optional(),
		field.Bool("verified").Default(false),
		field.Time("verified_at").Optional(),
	}
}

// Edges of the Business.
func (Business) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("socials", Social.Type),
		edge.To("manages", Manager.Type),
	}
}
