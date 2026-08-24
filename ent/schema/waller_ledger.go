package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Ledger holds the schema definition for the Ledger entity.
type Ledger struct {
	ent.Schema
}

func (Ledger) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Ledger.
func (Ledger) Fields() []ent.Field {
	return []ent.Field{
		field.String("transaction_id").NotEmpty().Unique(),
		field.String("reversal_transaction_id").Optional(),
		field.String("wallet_id").NotEmpty().Unique(),
		field.Int64("debit").Default(0),
		field.Int64("credit").Default(0),
		field.Int64("current_balance").Default(0),
		field.Int64("previous_balance").Default(0),
		field.Bool("is_reversal").
			Default(false),
	}
}

// Edges of the Ledger.
func (Ledger) Edges() []ent.Edge {
	return []ent.Edge{
		// Many LedgerEntries → One Transaction
		edge.From("transaction", Transaction.Type).
			Ref("ledger_entries").
			Field("transaction_id").
			Unique().
			Required(),

		// Many LedgerEntries → One Wallet
		edge.From("wallet", Wallet.Type).
			Ref("ledger_entries").
			Field("wallet_id").
			Unique().
			Required(),

		// Optional reversal relationship
		edge.From("reversal_transaction", Transaction.Type).
			Ref("reversal_ledger_entries").
			Field("reversal_transaction_id").
			Unique(),
	}
}
