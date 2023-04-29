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

type ExternalEdgeMixin struct {
	mixin.Schema

	config           *ExternalEdgeConfig
	fieldAnnotations []schema.Annotation
}

type ExternalEdgeConfig struct {
	EdgeName       string
	ConnectionName string
	FieldName      string
	AnyID          bool
	Immutable      bool
	Optional       bool
	SkipGQL        bool
}

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

// ExternalEdge is an ent.Mixin that adds the fields needed for external references
// to both the ent models as well as the generated gql schema
func ExternalEdge(opts ...ExternalEdgeOption) *ExternalEdgeMixin {
	m := &ExternalEdgeMixin{config: &ExternalEdgeConfig{}}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

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

func (m ExternalEdgeMixin) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields(m.config.FieldName),
	}
}
