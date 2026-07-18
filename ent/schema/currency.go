package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Currency holds the schema definition for the Currency entity.
type Currency struct {
	ent.Schema
}

func (Currency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Currency.
func (Currency) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("symbol").NotEmpty(),
		field.String("code").NotEmpty().Unique(),
		field.Bool("is_fiat").Default(true),
		field.Bool("active").Default(true),
		field.Int64("multipler"),
	}
}

// Edges of the Currency.
func (Currency) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("wallets", Wallet.Type),
	}
}
