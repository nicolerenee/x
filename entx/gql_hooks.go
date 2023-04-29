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

	addExternalEdges = func(g *gen.Graph, s *ast.Schema) error {
		for _, n := range g.Nodes {
			for _, f := range n.Fields {
				if an, ok := f.Annotations[AnnotationName]; ok {
					if m, ok := an.(map[string]interface{}); ok {
						if ed, ok := m["ExternalEdge"]; ok {
							if cfg, ok := ed.(map[string]interface{}); ok {

								edgeName := cfg["EdgeName"].(string)

								// Check if the external type is already defined and if not add it
								extDef, ok := s.Types[edgeName]
								if !ok {
									extDef = &ast.Definition{
										Kind:        ast.Object,
										Description: "",
										Name:        edgeName,
										Directives: []*ast.Directive{
											{
												Name: "key",
												Arguments: []*ast.Argument{
													{
														Name: "fields",
														Value: &ast.Value{
															Raw:  "id",
															Kind: ast.StringValue,
														},
													},
												},
											},
										},
										Fields: ast.FieldList{
											{
												Name: "id",
												Type: ast.NonNullNamedType("ID", nil),
												Directives: []*ast.Directive{
													{
														Name: "external",
													},
												},
											},
										},
									}
									s.AddTypes(extDef)
								}

								// Add the edge to the external type
								query, ok := s.Types["Query"]
								if !ok {
									return errors.New("failed to find query definition in schema")
								}

								connName := camel(snake(plural(n.Name)))
								if n := cfg["ConnectionName"].(string); n != "" {
									connName = n
								}

								qf := query.Fields.ForName(connName)
								if qf == nil {
									return errors.New("failed to find local edge in query definition")
								}

								extDef.Fields = append(extDef.Fields, qf)

								// Now add external edge to local edge
								localDef, ok := s.Types[n.Name]
								if !ok {
									return errors.New("unable to find local edge")
								}

								var localType *ast.Type
								if cfg["Optional"].(bool) {
									localType = ast.NamedType(edgeName, nil)
								} else {
									localType = ast.NonNullNamedType(edgeName, nil)
								}

								localDef.Fields = append(localDef.Fields, &ast.FieldDefinition{
									Name: camel(snake(edgeName)),
									Type: localType,
								})
							}
						}
					}
				}
				// an := Annotation{}
				// if name == an.Name() {
				// 	fmt.Printf("I'm adding external edges: %s: %s:\n\n", n.Name, name)
				// 	an = a.(Annotation)
				// 	if an.ExternalEdge == nil {
				// 		continue
				// 	}

			}
		}

		return nil
	}
)

// import string mutations from entc
var (
	_ entc.Extension = (*Extension)(nil)

	camel = gen.Funcs["camel"].(func(string) string)
	// pascal   = gen.Funcs["pascal"].(func(string) string)
	plural = gen.Funcs["plural"].(func(string) string)
	// singular = gen.Funcs["singular"].(func(string) string)
	snake = gen.Funcs["snake"].(func(string) string)
)
