package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// KYBMessage holds the schema definition for the KYBMessage entity.
type KYBMessage struct {
	ent.Schema
}

func (KYBMessage) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the KYBMessage.
func (KYBMessage) Fields() []ent.Field {
	return []ent.Field{
		field.String("message").NotEmpty(),
		field.Enum("status").Values("OPEN", "CLOSED", "RESOLVED").Default("OPEN"),
	}
}

// Edges of the KYBMessage.
func (KYBMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("business", Business.Type).
			Ref("kyb_messages").
			Required().
			Unique(),
	}
}
