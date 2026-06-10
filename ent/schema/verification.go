package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Verification holds the schema definition for the Verification entity.
type Verification struct {
	ent.Schema
}

func (Verification) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Verification.
func (Verification) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider").NotEmpty(),
		field.Enum("verification_type").Values(
			"BVN",
			"NIN",
			"PASSPORT",
			"VOTER_ID",
			"DRIVERS_LICENSE"),
		field.Enum("status").Values(
			"PENDING",
			"IN_PROGRESS",
			"VERIFIED",
			"FAILED",
			"REJECTED",
			"EXPIRED",
		).Default("IN_PROGRESS"),
		field.String("reference_id").NotEmpty(),
		field.String("provider_reference").NotEmpty(),
		field.JSON("metadata", map[string]interface{}{}),
		field.JSON("data", map[string]interface{}{}),
		field.Time("verified_at").Default(time.Now()).Optional(),
		field.Bool("is_obsolate").Default(false),
	}
}

// Edges of the Verification.
func (Verification) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type),
		edge.To("business", Business.Type),
	}
}
