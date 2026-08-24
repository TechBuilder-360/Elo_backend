package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type NubanDynamicAccount struct {
	ent.Schema
}

func (NubanDynamicAccount) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (NubanDynamicAccount) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider"),
		field.String("provider_reference"),
		field.String("account_number").NotEmpty(),
		field.String("account_name").NotEmpty(),
		field.String("bank_name").NotEmpty(),
		field.String("bank_code").NotEmpty(),
		field.Time("expiration").Nillable(),
		field.Enum("state").Values("OPEN", "SUSPENDED", "CLOSED").Default("OPEN"),
	}
}

func (NubanDynamicAccount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("wallet", Wallet.Type).
			Ref("nuban_dynamic_account").
			Unique().
			Required(),
	}
}
