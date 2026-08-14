package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Vault struct {
	ent.Schema
}

func (Vault) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Vault) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("type").
			Values("TREASURY").
			Default("TREASURY"),

		field.Enum("status").
			Values(
				"ACTIVE",
				"SUSPENDED",
				"CLOSED",
			).
			Default("ACTIVE"),

		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

func (Vault) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", LedgerOwner.Type).
			Ref("vaults").
			Unique().
			Required(),

		edge.To("wallets", Wallet.Type),
	}
}
