package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// StablecoinNetwork holds the schema definition for the StablecoinNetwork entity.
type StablecoinNetwork struct {
	ent.Schema
}

func (StablecoinNetwork) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StablecoinNetwork.
func (StablecoinNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").NotEmpty(),
		field.String("slug").NotEmpty().Unique(),
		field.String("logo_url").Nillable().Default(""),
		field.Bool("active").Default(true),
	}
}

// Edges of the StablecoinNetwork.
func (StablecoinNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("stablecoin_networks", StablecoinWallet.Type),
		edge.To("stablecoin_supported_networks", StablecoinSupportedNetwork.Type),
	}
}
