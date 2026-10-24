package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

func (Subscription) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("status").
			Values(
				"TRIALING",
				"ACTIVE",
				"PAST_DUE",
				"CANCELLED",
				"EXPIRED",
				"PAUSED",
			).
			Default("TRIALING"),

		// Snapshot of pricing at subscription time — protects against plan mutation/deletion
		field.Int64("locked_amount"),
		field.String("locked_currency").NotEmpty(),

		field.Time("current_period_start"),
		field.Time("current_period_end"),

		field.Time("trial_ends_at").Optional().Nillable(),
		field.Time("cancelled_at").Optional().Nillable(),
		field.Time("cancel_at_period_end").Optional().Nillable(), // graceful cancellation
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("plan", SubscriptionPlan.Type).
			Ref("subscriptions").
			Unique().
			Required(),

		edge.From("business", Business.Type).
			Ref("subscriptions").
			Unique().
			Required(),

		edge.To("invoices", Invoice.Type),
	}
}
