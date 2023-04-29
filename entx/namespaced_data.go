package entx

import (
	"encoding/json"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type NamespacedDataMixin struct {
	mixin.Schema
}

func (m NamespacedDataMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Text("namespace").
			NotEmpty().
			MinLen(5).
			MaxLen(64).
			Annotations(
				entgql.OrderField("NAMESPACE"),
			),
		field.JSON("data", json.RawMessage{}).
			Annotations(
				entgql.Type("JSON"),
				Annotation{IsNamespacedDataJSONField: true},
			),
	}
}
