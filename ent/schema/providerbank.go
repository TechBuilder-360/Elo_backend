package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type ProviderBank struct {
	ent.Schema
}

func (ProviderBank) Fields() []ent.Field {
	return []ent.Field{
		field.String("nuban_bank_id").
			NotEmpty(),

		field.String("provider_code").
			NotEmpty(),

		field.String("provider").
			NotEmpty(),

		field.Bool("active").
			Default(true),
	}
}

func (ProviderBank) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("bank", NubanBank.Type).
			Ref("provider_banks").
			Field("nuban_bank_id").
			Unique().
			Required(),
	}
}

func (ProviderBank) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider", "nuban_bank_id").
			Unique(),

		index.Fields("provider", "provider_code").
			Unique(),
	}
}
