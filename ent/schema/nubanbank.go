package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NubanBank holds the schema definition for the NubanBank entity.
type NubanBank struct {
	ent.Schema
}

func (NubanBank) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the NubanBank.
func (NubanBank) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty().Unique(),
		field.String("cbn_code").NotEmpty().Unique(),
		field.String("short_name").NotEmpty().Unique(),
		field.String("slug").NotEmpty().Unique(),
		field.Bool("active").Default(true),
	}
}

// Edges of the NubanBank.
func (NubanBank) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("provider_banks", ProviderBank.Type),
		edge.To("bank", NubanTransfer.Type),
		edge.To("beneficiary_bank", NubanBeneficiary.Type),
		edge.To("blacklist", NubanBlacklist.Type),
	}
}

func (NubanBank) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("cbn_code", "slug").
			Unique(),
	}
}
