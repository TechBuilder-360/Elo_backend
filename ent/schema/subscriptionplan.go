package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// SubscriptionPlan holds the schema definition for the SubscriptionPlan entity.
type SubscriptionPlan struct {
	ent.Schema
}

func (SubscriptionPlan) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the SubscriptionPlan.
func (SubscriptionPlan) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"), // e.g. "Pro Monthly"

		field.String("slug").Unique(), // e.g. "pro-monthly", stable identifier for code references

		field.Text("description").Optional(),

		field.Enum("billing_interval").
			Values(
				"DAILY",
				"WEEKLY",
				"MONTHLY",
				"QUARTERLY",
				"YEARLY",
			),

		field.Int("interval_count").Default(1), // e.g. every 3 months = MONTHLY + interval_count=3

		field.Int64("amount"), // smallest unit — recurring charge per interval

		field.String("currency").NotEmpty(),

		field.Int("trial_period_days").Default(0),

		field.JSON("features", []string{}).Optional(), // or a structured entitlements map

		field.Int("max_seats").Optional().Nillable(), // if plans are seat-based

		field.Bool("is_active").Default(true), // whether new subscriptions can be created on this plan

		field.Bool("is_public").Default(true), // visible in pricing page vs. legacy/custom plan
	}
}

// Edges of the Subscription.
func (SubscriptionPlan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("subscriptions", Subscription.Type),
	}
}
