package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Manager holds the schema definition for the Manager entity.
type Manager struct {
	ent.Schema
}

// Fields of the Manager.
func (Manager) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.String("position").NotEmpty(),
		field.Bool("diasbled").Default(false),
		field.String("disable_reason").Nillable(),
		field.Time("disabled_at").Nillable(),
	}
}

// Edges of the Manager.
func (Manager) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("admins", Business.Type).
			Ref("manager").
			Unique().
			Required(),
		edge.From("managers", User.Type).
			Ref("manager").
			Unique().
			Required(),
	}
}

// func (Manager) Indexes() []ent.Index {
// 	return []ent.Index{
// 		index.
// 			Fields("user_id", "business_id").
// 			Unique(), // Enforces composite uniqueness
// 	}
// }
