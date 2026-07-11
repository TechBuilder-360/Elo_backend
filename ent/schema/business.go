package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/Toflex/directory_v2/pkg/util"
)

// Business holds the schema definition for the Business entity.
type Business struct {
	ent.Schema
}

func (Business) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Business.
func (Business) Fields() []ent.Field {
	return []ent.Field{
		field.String("category").Default("others"),
		field.String("name").Unique().NotEmpty(),
		field.String("about").NotEmpty(),
		field.String("logo").Optional(),
		field.String("cover_image").Optional(),
		field.String("registered_by").Nillable().Optional(),
		field.String("country_of_incorporation").Nillable().Optional(),
		field.String("date_of_incorporation").Nillable().Optional(),
		field.String("registration_number").Nillable().Optional(),
		field.String("email").NotEmpty().
			Validate(util.ValidateEmail),
		field.String("website").Optional().
			Validate(util.ValidateURL),
		field.Bool("on_site").Default(false),
		field.Bool("active").Default(false),
		field.Bool("live").Default(false),
		field.Bool("disabled").Default(true),
		field.Time("disabled_at").Default(time.Now),
		field.String("disable_reason").Optional(),
		field.Enum("verification_status").Values(
			"UNVERIFIED",
			"IN_PROGRESS",
			"VERIFIED",
			"REJECTED",
		).Default("UNVERIFIED"),
		field.Bool("verified").Default(false),
		field.Time("verified_at").Optional(),
	}
}

// Edges of the Business.
func (Business) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("socials", Social.Type),
		edge.To("services", BusinessServices.Type),
		edge.To("manages", Manager.Type),
		edge.From("registered_by_user", User.Type).
			Ref("registered_businesses").
			Field("registered_by").
			Unique(),
		edge.From("verifications", Verification.Type).
			Ref("business"),
		edge.From("request_verifications", RequestVerification.Type).
			Ref("business"),
		edge.To("business_documents", BusinessDocument.Type),
		edge.To("locations", BusinessLocation.Type),
		edge.To("kyb_messages", KYBMessage.Type),
	}
}

func (Business) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("registration_number", "name").Unique(),
	}
}
