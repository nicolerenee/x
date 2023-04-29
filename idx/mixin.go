package idx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/vektah/gqlparser/v2/ast"
)

// AnnotationName is the value of the annotation when read during ent compilation
const AnnotationName string = "PrefixIDPrefix"

// PrimaryKeyMixin creates a Mixin that provides the primary key as a PrefixedID.
func PrimaryKeyMixin(t string) *Mixin {
	return &Mixin{
		prefix:    t,
		isIDField: true,
	}
}

// Mixin defines an ent Mixin that captures a field as a NamespaceUUID.
type Mixin struct {
	mixin.Schema
	prefix    string
	isIDField bool
}

// Fields provides the id field.
func (m Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(PrefixedID("")).
			DefaultFunc(func() PrefixedID { return MustNewID(m.prefix) }).
			Unique().
			Immutable(),
	}
}

// Annotation captures the id prefix for a type.
type Annotation struct {
	Prefix string
}

// Name implements the ent Annotation interface.
func (a Annotation) Name() string {
	return AnnotationName
}

// Annotations returns the annotations for a Mixin instance.
func (m Mixin) Annotations() []schema.Annotation {
	ans := []schema.Annotation{}

	if m.isIDField {
		ans = append(ans,
			Annotation{Prefix: m.prefix},
			entgql.Directives(KeyDirective("id")),
		)
	}

	return ans
}

// KeyDirective returns an entgql.Directive for setting the @key field on a gql
// type
func KeyDirective(fields string) entgql.Directive {
	var args []*ast.Argument
	if fields != "" {
		args = append(args, &ast.Argument{
			Name: "fields",
			Value: &ast.Value{
				Raw:  fields,
				Kind: ast.StringValue,
			},
		})
	}

	return entgql.NewDirective("key", args...)
}
