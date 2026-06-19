package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Toflex/directory_v2/pkg/util"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").NotEmpty(),
		field.String("last_name").NotEmpty(),
		field.String("password").NotEmpty(),
		field.String("middle_name").Optional(),
		field.String("display_name").Optional(),
		field.String("email_address").NotEmpty().Unique().Validate(util.ValidateEmail),
		field.Bool("email_verified").Default(false),
		field.Time("email_verified_at").Optional(),
		field.String("phone_number").Optional().Nillable().Validate(util.IsValidPhoneNumber),
		field.String("avatar").Optional().Nillable(),
		field.Bool("disabled").Default(false),
		field.String("disable_reason").Nillable().Optional(),
		field.Bool("verified").Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manages", Manager.Type),
		edge.To("user_documents", UserDocument.Type),
		edge.From("verifications", Verification.Type).
			Ref("user"),
		edge.From("request_verifications", RequestVerification.Type).
			Ref("user"),
	}
}
