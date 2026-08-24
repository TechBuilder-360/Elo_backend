package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StablecoinDeposit holds the schema definition for the StablecoinDeposit entity.
type StablecoinDeposit struct {
	ent.Schema
}

func (StablecoinDeposit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StablecoinDeposit.
func (StablecoinDeposit) Fields() []ent.Field {
	return []ent.Field{
		field.String("coin").NotEmpty(),
		field.String("network").NotEmpty(),
		field.String("address").NotEmpty(),
		field.String("transaction_id").NotEmpty(),
		field.String("stablecoin_wallet_id").NotEmpty(),
		field.String("provider_reference").NotEmpty(),
		field.String("provider").NotEmpty(),
	}
}

// Edges of the StablecoinDeposit.
func (StablecoinDeposit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("stablecoin_deposit").
			Field("transaction_id").
			Unique().
			Required(),
		edge.From("stablecoin_wallet", StablecoinWallet.Type).
			Ref("stablecoin_deposits").
			Field("stablecoin_wallet_id").
			Unique().
			Required(),
	}
}
