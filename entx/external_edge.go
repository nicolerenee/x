package entx

import (
	"fmt"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
	"github.com/gobeam/stringy"

	"go.infratographer.com/x/idx"
)

// ExternalEdgeMixin defines an ent Mixin that allows you to define an Edge that
// is external to the current system.
type ExternalEdgeMixin struct {
	mixin.Schema

	config           *ExternalEdgeConfig
	fieldAnnotations []schema.Annotation
}

// ExternalEdgeConfig provides all the final config options used for an External
// Edge. These are exposed via the annotation on the field for use in hooks and
// templates.
type ExternalEdgeConfig struct {
	EdgeName       string
	ConnectionName string
	FieldName      string
	AnyID          bool
	Immutable      bool
	Optional       bool
	SkipGQL        bool
}

// ExternalEdgeOption provides the ability to customize the behavior of ExternalEdges
type ExternalEdgeOption func(*ExternalEdgeMixin)

// WithExternalEdgeType allows you to define the gql type of the external edge.
// The type will be used to generate a default field name.
func WithExternalEdgeType(name string) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		if name != stringy.New(name).CamelCase() {
			fmt.Printf(
				"\nWARNING: The external edge does not appear to be formatted correctly, please verify this value before continuing. Got: %s, Expected: %s\n",
				name,
				stringy.New(name).CamelCase(),
			)
		}

		m.config.EdgeName = name

		if m.config.FieldName == "" {
			s := stringy.New(name).SnakeCase().ToLower()
			m.config.FieldName = fmt.Sprintf("%s_id", s)
		}
	}
}

// ExternalEdgeOptional will set the field to be an optional field. By default
// an external edge is required
func ExternalEdgeOptional(v bool) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		m.config.Optional = v
	}
}

// ExternalEdgeImmutable will set the field to be an immutable field.
func ExternalEdgeImmutable(v bool) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		m.config.Immutable = v
	}
}

// SkipExternalEdgeGQL will skip the creation of gql schema to support this external
// edge.
func SkipExternalEdgeGQL(v bool) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		m.config.SkipGQL = v
	}
}

// WithExternalEdgeFieldAnnotations allows you to add additional annotations
// to the resulting field used to store the ID of the external edge
func WithExternalEdgeFieldAnnotations(a entgql.Annotation) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		m.fieldAnnotations = append(m.fieldAnnotations, a)
	}
}

// ExternalEdgeAllowAnyID will set the field to allow any id values. By default
// the infratographer idx.PrefixedID values are required.
func ExternalEdgeAllowAnyID(v bool) ExternalEdgeOption {
	return func(m *ExternalEdgeMixin) {
		m.config.AnyID = v
	}
}

// ExternalEdge returns an ent.Mixin that adds the fields needed for external edge
// references to both the ent models as well as the generated gql schema. This
// makes managing federated gql easier since the schema extensions for the
// remote service types are added automatically to your schema.
func ExternalEdge(opts ...ExternalEdgeOption) *ExternalEdgeMixin {
	m := &ExternalEdgeMixin{config: &ExternalEdgeConfig{}}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Fields provides the id field of the ExternalEdge
func (m ExternalEdgeMixin) Fields() []ent.Field {
	f := field.Text(m.config.FieldName).Annotations(
		entgql.Type("ID"),
		Annotation{ExternalEdge: m.config},
	)

	if m.config.Immutable {
		f.Immutable().Annotations(entgql.Skip(entgql.SkipMutationUpdateInput))
	}

	if !m.config.AnyID {
		f.MinLen(idx.TotalLength).MaxLen(idx.TotalLength).GoType(idx.PrefixedID(""))
	}

	if !m.config.Optional {
		f.Optional()
	} else {
		f.NotEmpty()
	}

	// provided annotations are applied last so that they can override other annotations
	f.Annotations(m.fieldAnnotations...)

	return []ent.Field{f}
}

// Indexes provides the index of the field representing an ExternalEdge
func (m ExternalEdgeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(m.config.FieldName),
	}
}
