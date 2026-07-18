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
			Values("TREASURY", "HOLDING").
			Default("TREASURY"),
		field.Int64("available_balance").Default(0),
		field.Int64("ledger_balance").Default(0),
		field.Int64("holding_balance").Default(0),
		field.Enum("owner").
			Values("USER", "BUSINESS").
			Default("USER"),
		field.String("identifier").
			Comment("These field avoids the use of fk, by using naming convension 'user-user_id' of 'business-business_id'"),
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
	}
}

func (Wallet) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("identifier", "currency_id", "type").Unique(),
	}
}
