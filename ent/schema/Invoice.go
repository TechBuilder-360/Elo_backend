package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Invoice holds the schema definition for the Invoice entity.
type Invoice struct {
	ent.Schema
}

func (Invoice) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (Invoice) Fields() []ent.Field {
	return []ent.Field{
		field.String("reference").Unique(),

		field.Enum("status").
			Values("PENDING", "PAID", "FAILED", "VOID", "REFUNDED"),

		field.Int64("amount"), // gross, before fees
		field.String("currency").NotEmpty(),

		field.Time("period_start"),
		field.Time("period_end"),

		field.Time("due_at"),
		field.Time("paid_at").Optional().Nillable(),
	}
}

func (Invoice) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subscription", Subscription.Type).
			Ref("invoices").
			Unique().
			Required(),

		edge.From("transaction", Transaction.Type).
			Ref("invoice").
			Unique().
			Required(),
	}
}
