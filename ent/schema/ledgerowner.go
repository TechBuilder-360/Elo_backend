package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type LedgerOwner struct {
	ent.Schema
}

func (LedgerOwner) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (LedgerOwner) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values(
				"BUSINESS",
				"USER",
			),

		field.Enum("status").
			Values(
				"ACTIVE",
				"SUSPENDED",
				"CLOSED",
			).
			Default("ACTIVE"),
	}
}

func (LedgerOwner) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vaults", Vault.Type),

		edge.From("business", Business.Type).
			Ref("owner").
			Unique(),

		edge.From("user", User.Type).
			Ref("owner").
			Unique(),
	}
}
