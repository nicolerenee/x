package entx

import (
	"entgo.io/contrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

// GraphQLKeyDirective returns an entgql.Directive for setting the @key field on
// a graphql type
func GraphKeyDirective(key string) entgql.Annotation {
	return entgql.Directives(keyDirective("id"))
}

func keyDirective(fields string) entgql.Directive {
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
