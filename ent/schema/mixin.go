package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Int64IdMixin defines a reusable ID field of type int64
type Int64IdMixin struct {
	mixin.Schema
}

// Fields are the fields of the Int64IdMixin
func (Int64IdMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
	}
}
