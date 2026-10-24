package schema

import (
	"errors"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FeeRule holds the schema definition for the FeeRule entity.
type FeeRule struct {
	ent.Schema
}

func (FeeRule) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the FeeRule.
func (FeeRule) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"), // e.g. "Standard Withdrawal VAT"

		field.String("service_id"),

		field.Enum("fee_type"). // matches your TransactionFee.type
					Values(
				"VAT",
				"SERVICE_FEE",
				"NETWORK_FEE",
				"PROCESSING_FEE",
				"WITHHOLDING_TAX",
			),

		field.Enum("calculation_method").
			Values(
				"FLAT",
				"PERCENTAGE",
				"TIERED",
			),

		// Used when calculation_method == FLAT
		field.Int64("flat_amount").Optional().Default(0),

		// Used when calculation_method == PERCENTAGE
		field.Float("percentage_rate").Validate(func(f float64) error {
			if f > 100 {
				return errors.New("percentage_rate cannot exceed 100")
			}
			if f < 0 {
				return errors.New("percentage_rate cannot be negative")
			}

			return nil
		}).Optional().Default(0),

		field.String("currency").NotEmpty(),

		field.Bool("is_active").Default(true),

		field.Time("effective_from").Default(time.Now),
		field.Time("effective_to").Optional().Nillable(),

		// Eligibility window on the transaction amount (smallest unit, both inclusive).
		// nil min = no lower bound, nil max = no upper bound.
		field.Int64("min_value").Optional().Nillable().NonNegative(),
		field.Int64("max_value").Optional().Nillable().NonNegative(),

		// Optional guardrails on the computed fee (e.g. "1.5%, min 50, capped at 2000").
		field.Int64("min_fee").Optional().Nillable().NonNegative(),
		field.Int64("max_fee").Optional().Nillable().NonNegative(),
	}
}

// Edges of the FeeRule.
func (FeeRule) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("service", Service.Type).
			Ref("fee_rules").
			Field("service_id").
			Unique().
			Required(),

		edge.To("tiers", FeeTier.Type), // only populated when calculation_method == TIERED
	}
}
