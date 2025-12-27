package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Toflex/directory_v2/pkg/utils"
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
		field.String("middle_name").Optional(),
		field.String("display_name").Optional(),
		field.String("email_address").NotEmpty().Unique().Validate(utils.ValidateEmail),
		field.Bool("email_verified").Default(false),
		field.Time("email_verified_at").Optional(),
		field.String("phone_number").Optional().Validate(utils.IsValidPhoneNumber),
		field.String("avatar").Optional(),
		field.Bool("disabled").Default(false),
		field.Int8("tier").Default(0),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manages", Manager.Type),
	}
}
