package schema

import (
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FeeTier holds the schema definition for the FeeTier entity.
type FeeTier struct {
	ent.Schema
}

func (FeeTier) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the FeeTier.
func (FeeTier) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("min_amount").Default(0),            // inclusive, smallest unit
		field.Int64("max_amount").Optional().Nillable(), // inclusive; nil = unbounded (last tier)

		field.Enum("charge_type").
			Values("FLAT", "PERCENTAGE"),

		field.Int64("flat_amount").Optional().Default(0),
		field.Float("percentage_rate").Optional().Validate(func(f float64) error {
			if f > 100 {
				return errors.New("percentage_rate cannot exceed 100")
			}
			if f < 0 {
				return errors.New("percentage_rate cannot be negative")
			}

			return nil
		}).Default(0),

		field.Int("priority").Default(0), // ordering, in case ranges are evaluated sequentially
	}
}

// Edges of the FeeTier.
func (FeeTier) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fee_rule", FeeRule.Type).
			Ref("tiers").
			Unique().
			Required(),
	}
}
