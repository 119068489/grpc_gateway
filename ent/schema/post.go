package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Post struct {
	ent.Schema
}

// Fields of the Post.
func (Post) Fields() []ent.Field {
	return []ent.Field{
		field.Int("author_id").
			// StorageKey("post_author").
			Optional().
			Nillable(),
	}
}

// Edges of the Post.
func (Post) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).
			// Bind the "author_id" field to this edge.
			Field("author_id").
			StorageKey(edge.Column("post_author")).
			Unique(),
		// edge.From("author", User.Type).
		// 	Field("author_id").
		// 	StorageKey(edge.Column("post_author")).Unique(),
	}
}
