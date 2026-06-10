package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Toflex/directory_v2/pkg/util"
)

// RequestVerification holds the schema definition for the RequestVerification entity.
type RequestVerification struct {
	ent.Schema
}

func (RequestVerification) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the RequestVerification.
func (RequestVerification) Fields() []ent.Field {
	return []ent.Field{
		field.String("reference_id").Unique().NotEmpty(),
		field.String("verification_type").NotEmpty(),
		field.String("provider").NotEmpty(),
		field.String("link").NotEmpty().Validate(func(s string) error {
			return util.ValidateURL(s)
		}),
		field.String("provider_link").NotEmpty().Validate(func(s string) error {
			return util.ValidateURL(s)
		}),
		field.Enum("status").Values(
			"PENDING",
			"IN_PROGRESS",
			"VERIFIED",
			"FAILED",
			"REJECTED",
			"EXPIRED",
		).Default("PENDING"),
	}
}

// Edges of the RequestVerification.
func (RequestVerification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("business", Business.Type),
	}
}
