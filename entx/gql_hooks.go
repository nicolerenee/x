package entx

import (
	"errors"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

// Skipping err113 linting since these errors are returned during generation and not runtime
//
//nolint:goerr113
var (
	removeNodeGoModel = func(g *gen.Graph, s *ast.Schema) error {
		n, ok := s.Types["Node"]
		if !ok {
			return errors.New("failed to find node interface in schema")
		}

		dirs := ast.DirectiveList{}

		for _, d := range n.Directives {
			switch d.Name {
			case "goModel":
				continue
			default:
				dirs = append(dirs, d)
			}
		}
		n.Directives = dirs

		return nil
	}

	removeNodeQueries = func(g *gen.Graph, s *ast.Schema) error {
		q, ok := s.Types["Query"]
		if !ok {
			return errors.New("failed to find query definition in schema")
		}

		fields := ast.FieldList{}

		for _, f := range q.Fields {
			switch f.Name {
			case "node":
			case "nodes":
				continue
			default:
				fields = append(fields, f)
			}
		}
		q.Fields = fields

		return nil
	}

	addJSONScalar = func(g *gen.Graph, s *ast.Schema) error {
		s.Types["JSON"] = &ast.Definition{
			Kind:        ast.Scalar,
			Description: "A valid JSON string.",
			Name:        "JSON",
		}
		return nil
	}
)

// import string mutations from entc
var (
	_ entc.Extension = (*Extension)(nil)

	camel  = gen.Funcs["camel"].(func(string) string)
	pascal = gen.Funcs["pascal"].(func(string) string)
	plural = gen.Funcs["plural"].(func(string) string)
	// singular = gen.Funcs["singular"].(func(string) string)
	snake = gen.Funcs["snake"].(func(string) string)
)
