package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Wallet holds the schema definition for the Wallet entity.
type Wallet struct {
	ent.Schema
}

func (Wallet) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Wallet.
func (Wallet) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values("FIAT", "CRYPTO").
			Default("FIAT"),
		field.Int64("available_balance").Default(0),
		field.Int64("ledger_balance").Default(0),
		field.Int64("holding_balance").Default(0),
		field.String("currency_id").
			NotEmpty(),
		field.Bool("active").
			Default(true),
	}
}

// Edges of the Wallet.
func (Wallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("currency", Currency.Type).
			Ref("wallets").
			Field("currency_id").
			Required().
			Unique(),

		edge.From("vault", Vault.Type).
			Ref("wallets").
			Unique().
			Required(),

		edge.To("nuban_static_account", NubanStaticAccount.Type),
	}
}

func (Wallet) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("currency_id").
			Edges("vault").
			Unique(),
	}
}
