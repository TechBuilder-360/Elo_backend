package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/Toflex/directory_v2/pkg/utils"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("first_name").NotEmpty(),
		field.String("last_name").NotEmpty(),
		field.String("middle_name").Nillable(),
		field.String("display_name").Nillable(),
		field.String("email_address").NotEmpty().Unique().Validate(utils.ValidateEmail),
		field.String("phone_number").Nillable().Validate(utils.IsValidPhoneNumber),
		field.String("avatar").Nillable(),
		field.Bool("disabled").Default(false),
		field.String("identification_number").Nillable(),
		field.Int8("tier").Default(0),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("manager", Manager.Type),
	}
}
