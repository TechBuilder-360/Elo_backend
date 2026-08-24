package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Transaction holds the schema definition for the Transaction entity.
type Transaction struct {
	ent.Schema
}

func (Transaction) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Transaction.
func (Transaction) Fields() []ent.Field {
	return []ent.Field{
		field.String("reference").
			Unique(),
		field.String("external_reference").Nillable().
			Unique(),
		field.Enum("type").
			Values(
				"DEPOSIT",
				"WITHDRAWAL",
				"TRANSFER",
				"PAYMENT",
				"REFUND",
				"REVERSAL",
				"FX",
				"FEE",
				"ADJUSTMENT",
			),
		field.Enum("channel").
			Values(
				"CRYPTO",
				"BANK_TRANSFER",
				"INTERNAL",
			),
		field.Enum("status").
			Values(
				"PENDING",
				"PROCESSING",
				"COMPLETED",
				"FAILED",
				"REVERSED",
				"CANCELLED",
			).
			Default("PENDING"),
		field.String("currency").NotEmpty(),
		field.String("summary").Optional(),
		field.String("provider").
			Optional(),
		field.Int64("amount").Default(0),
		field.Int64("fee").Default(0),
	}
}

// Edges of the Transaction.
func (Transaction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("nuban_deposit", NubanDeposit.Type).
			Unique(),
		edge.To("nuban_transfer", NubanTransfer.Type).
			Unique(),
		edge.To("wallet", Wallet.Type).
			Unique().
			Required(),
		edge.To("stablecoin_deposit", StablecoinDeposit.Type).
			Unique(),
		edge.To("stablecoin_withdrawal", StablecoinWithdrawal.Type).
			Unique(),

		// LEDGER
		edge.To("ledger_entries", Ledger.Type),
		edge.To("reversal_ledger_entries", Ledger.Type),
	}
}
