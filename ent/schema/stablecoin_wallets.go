package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StablecoinWallet holds the schema definition for the StablecoinWallet entity.
type StablecoinWallet struct {
	ent.Schema
}

func (StablecoinWallet) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StablecoinWallet.
func (StablecoinWallet) Fields() []ent.Field {
	return []ent.Field{
		field.String("coin").NotEmpty(),
		field.String("network").NotEmpty(),
		field.String("network_id").NotEmpty(),
		field.String("coin_id").NotEmpty(),
		field.String("address").NotEmpty(),
		field.String("provider_reference").NotEmpty(),
		field.String("provider").NotEmpty(),
		field.Bool("disabled").Default(false),
	}
}

// Edges of the StablecoinWallet.
func (StablecoinWallet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).
			Ref("stablecoin_wallets").
			Unique().
			Required(),

		edge.From("currency", Currency.Type).
			Ref("stablecoin_currencies").
			Field("coin_id").
			Required().
			Unique(),

		edge.From("stablecoin_networks", StablecoinNetwork.Type).
			Ref("stablecoin_networks").
			Field("network_id").
			Required().
			Unique(),

		edge.To("stablecoin_deposits", StablecoinDeposit.Type),
	}
}
