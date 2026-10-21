package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// NubanBeneficiary stores beneficiary details for NUBAN transfers.
type NubanBeneficiary struct {
	ent.Schema
}

func (NubanBeneficiary) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

func (NubanBeneficiary) Fields() []ent.Field {
	return []ent.Field{
		field.String("account_name").NotEmpty(),
		field.String("account_number").NotEmpty(),
		field.String("bank_name").NotEmpty(),
		field.String("nuban_bank_id").NotEmpty(),
		field.String("owner_id").NotEmpty(),
		field.String("nickname").Optional().Nillable(),
		field.Bool("disabled").Default(false),
	}
}

func (NubanBeneficiary) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("nuban_beneficiary", LedgerOwner.Type).
			Ref("nuban_beneficiaries").
			Field("owner_id").
			Unique().
			Required(),
		edge.To("nuban_beneficiaries", NubanTransfer.Type),
		edge.From("nuban_bank", NubanBank.Type).
			Ref("beneficiary_bank").
			Field("nuban_bank_id").
			Unique().
			Required(),
	}
}

func (NubanBeneficiary) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id", "account_number", "nuban_bank_id").Unique(),
		index.Fields("owner_id"),
	}
}
