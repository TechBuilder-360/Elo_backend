package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// StablecoinSupportedNetwork holds the schema definition for the StablecoinSupportedNetwork entity.
type StablecoinSupportedNetwork struct {
	ent.Schema
}

func (StablecoinSupportedNetwork) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the StablecoinSupportedNetwork.
func (StablecoinSupportedNetwork) Fields() []ent.Field {
	return []ent.Field{
		field.String("network_id").NotEmpty(),
		field.String("coin_id").NotEmpty(),
		field.Bool("can_send").Default(false),
		field.Bool("can_receive").Default(true),
		field.Bool("active").Default(true),
	}
}

// Edges of the StablecoinSupportedNetwork.
func (StablecoinSupportedNetwork) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("currency", Currency.Type).
			Ref("stablecoin_supported_networks").
			Field("coin_id").
			Required().
			Unique(),
		edge.From("stablecoin_network", StablecoinNetwork.Type).
			Ref("stablecoin_supported_networks").
			Field("network_id").
			Required().
			Unique(),
	}
}

func (StablecoinSupportedNetwork) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("network_id", "coin_id").
			Unique(),
	}
}
