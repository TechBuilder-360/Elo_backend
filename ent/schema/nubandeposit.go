package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// NubanDeposit holds the schema definition for the NubanDeposit entity.
type NubanDeposit struct {
	ent.Schema
}

func (NubanDeposit) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the NubanDeposit.
func (NubanDeposit) Fields() []ent.Field {
	return []ent.Field{
		field.String("recipient_account_name").NotEmpty(),
		field.String("recipient_account_number").NotEmpty(),
		field.String("sender_account_name").NotEmpty(),
		field.String("sender_account_number").NotEmpty(),
		field.String("sender_bank_name").NotEmpty(),
		field.String("sender_bank_code").NotEmpty(),
		field.String("narration").NotEmpty(),
		field.Int64("amount").Default(0),
		field.String("transaction_id").NotEmpty(),
		field.String("session_id").NotEmpty(),
		field.String("provider_reference").NotEmpty(),
		field.String("provider").NotEmpty(),
		field.Enum("type").
			Values("DYNAMIC", "STATIC").
			Default("STATIC"),
	}
}

// Edges of the NubanDeposit.
func (NubanDeposit) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("transaction", Transaction.Type).
			Ref("nuban_deposit").
			Field("transaction_id").
			Unique().
			Required(),
	}
}
