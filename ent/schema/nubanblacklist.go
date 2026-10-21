package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NubanBlacklist stores beneficiary details for NUBAN transfers.
type NubanBlacklist struct {
	ent.Schema
}

func (NubanBlacklist) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (NubanBlacklist) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_name").NotEmpty(),
		field.String("account_number").NotEmpty(),
		field.String("bank_name").NotEmpty(),
		field.String("nuban_bank_id").NotEmpty(),
		field.String("reason").NotEmpty(),
		field.Bool("disabled").Default(true),
	}
}

func (NubanBlacklist) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("nuban_bank", NubanBank.Type).
			Ref("blacklist").
			Field("nuban_bank_id").
			Unique().
			Required(),
	}
}

func (NubanBlacklist) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("account_number", "nuban_bank_id").Unique(),
	}
}
