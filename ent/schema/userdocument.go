package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserDocument holds the schema definition for the UserDocument entity.
type UserDocument struct {
	ent.Schema
}

// Fields of the UserDocument.
func (UserDocument) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").NotEmpty(),
		field.String("description").NotEmpty(),
		field.String("url").NotEmpty(),
	}
}

// Edges of the UserDocument.
func (UserDocument) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user_document", Business.Type).
			Ref("user_documents").
			Required().
			Unique(),
	}
}
