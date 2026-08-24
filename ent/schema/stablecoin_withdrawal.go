package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StablecoinWithdrawal holds the schema definition for the StablecoinWithdrawal entity.
type StablecoinWithdrawal struct {
	ent.Schema
}

func (StablecoinWithdrawal) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StablecoinWithdrawal.
func (StablecoinWithdrawal) Fields() []ent.Field {
	return []ent.Field{
		field.String("coin").NotEmpty(),
		field.String("network").NotEmpty(),
		field.String("destination_address").NotEmpty(),
		field.Int64("amount").Default(0),
		field.String("transaction_id").NotEmpty(),
		field.String("provider_reference").Optional(),
		field.String("provider").NotEmpty(),
	}
}

// Edges of the StablecoinWithdrawal.
func (StablecoinWithdrawal) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("stablecoin_withdrawal").
			Field("transaction_id").
			Unique().
			Required(),
	}
}
