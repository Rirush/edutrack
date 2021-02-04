package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Lecture holds the schema definition for the Lecture entity.
type Lecture struct {
	ent.Schema
}

// Fields of the Lecture.
func (Lecture) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("title"),
		field.String("body"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").UpdateDefault(time.Now).Default(time.Now),
	}
}

// Edges of the Lecture.
func (Lecture) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("subject", Subject.Type).
			Required().
			Unique().
			Ref("lectures"),
	}
}
