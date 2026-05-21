package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Service holds the schema definition for the Service entity.
type Service struct {
	ent.Schema
}

type FeeType string

const (
	TierFeeType       FeeType = "TIER"
	PercentageFeeType FeeType = "PERCENTAGE"
	FlatFeeType       FeeType = "FLAT"
)

type Tier struct {
	From  uint
	To    uint
	Type  FeeType
	Value float64
}

type Fee struct {
	Type  FeeType
	Value float64
	Min   int
	Max   int
	Tiers []Tier
}

func (Service) Mixin() []ent.Mixin {
	return []ent.Mixin{
		BaseMixin{},
	}
}

// Fields of the Service.
func (Service) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique(),
		field.String("identifier").
			NotEmpty().
			Unique(),
		field.String("provider").
			NotEmpty().
			Unique(),
		field.Bool("require_subscription").
			Default(false),
		field.Bool("active").
			Default(true),
		field.Int("min").
			Default(0),
		field.Int("max").
			Default(0),
		field.JSON("fee", &Fee{}),
	}
}

// Edges of the Service.
func (Service) Edges() []ent.Edge {
	return nil
}

// Indexes of the Service.
func (Service) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("identifier").Unique(),
	}
}
