package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// TransactionFee holds the schema definition for the Transaction entity.
type TransactionFee struct {
	ent.Schema
}

func (TransactionFee) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the TransactionFee.
func (TransactionFee) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values(
				"VAT",
				"SERVICE_FEE",
				"NETWORK_FEE",
				"PROCESSING_FEE",
				"WITHHOLDING_TAX",
				"FX_MARKUP",
				"PLATFORM_FEE",
			),

		field.Int64("amount"), // smallest unit, same currency as parent txn (or its own — see below)

		field.String("currency").NotEmpty(),

		// Optional: if this fee was computed as a % of something (e.g. VAT = 7.5% of amount)
		field.Float("rate").Optional().Nillable(),

		field.Int64("base_amount").Optional().Nillable(), // the amount the rate was applied to

		field.String("service_identifier").Optional(), // provider-specific breakdown
	}
}

// Edges of the TransactionFee.
func (TransactionFee) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("fees").
			Unique().
			Required(),
	}
}
