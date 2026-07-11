package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Manager holds the schema definition for the Manager entity.
type Manager struct {
	ent.Schema
}

func (Manager) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Manager.
func (Manager) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("business_id"),
		field.Bool("is_owner").Default(false),
		field.String("role_id").NotEmpty(),
		field.Bool("disabled").Default(false),
		field.String("disable_reason").Optional(),
		field.Time("disabled_at").Optional(),
	}
}

// Edges of the Manager.
func (Manager) Edges() []ent.Edge {

	return []ent.Edge{
		edge.From("business", Business.Type).
			Ref("manages").
			Field("business_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("manages").
			Field("user_id").
			Required().
			Unique(),
		edge.From("role", Role.Type).
			Ref("managers").
			Field("role_id").
			Required().
			Unique(),
	}
}

func (Manager) Indexes() []ent.Index {
	return []ent.Index{
		index.
			Fields("user_id", "business_id", "role_id").
			Unique(), // Enforces composite uniqueness
	}
}
