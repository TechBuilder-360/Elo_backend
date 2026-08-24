package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NubanTransfer holds the schema definition for the NubanTransfer entity.
type NubanTransfer struct {
	ent.Schema
}

func (NubanTransfer) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the NubanTransfer.
func (NubanTransfer) Fields() []ent.Field {
	return []ent.Field{
		field.String("recipient_account_name").NotEmpty(),
		field.String("recipient_account_number").NotEmpty(),
		field.String("sender_account_name").NotEmpty(),
		field.String("sender_account_number").NotEmpty(),
		field.String("sender_bank_name").NotEmpty(),
		field.String("sender_bank_code").NotEmpty(),
		field.Int64("amount").Default(0),
		field.String("session_id").NotEmpty(),
		field.String("provider_reference").NotEmpty(),
		field.String("provider").NotEmpty(),
	}
}

// Edges of the NubanTransfer.
func (NubanTransfer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("nuban_transfer").
			Unique().
			Required(),
	}
}
