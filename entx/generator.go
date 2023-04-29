package entx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

// Extension is an implementation of entc.Extension that adds all the templates
// that entx needs.
type Extension struct {
	entc.DefaultExtension

	templates []*gen.Template

	gqlSchemaHooks []entgql.SchemaHook
}

type ExtensionOption func(*Extension) error

// WithFederation adds support for graphql federation by adding the Entity interface
// to all types, as well as removing the node() and nodes() query calls.
func WithFederation() ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = append(ex.templates, FederationTemplate)
		ex.gqlSchemaHooks = append(ex.gqlSchemaHooks, removeNodeGoModel, removeNodeQueries)
		return nil
	}
}

// WithJSONScalar adds the JSON scalar definition
func WithJSONScalar() ExtensionOption {
	return func(ex *Extension) error {
		ex.gqlSchemaHooks = append(ex.gqlSchemaHooks, addJSONScalar)
		return nil
	}
}

// WithGQLExternalEdges adds any external edges to the gql schema
func WithGQLExternalEdges() ExtensionOption {
	return func(ex *Extension) error {
		ex.gqlSchemaHooks = append(ex.gqlSchemaHooks, addExternalEdges)
		return nil
	}
}

func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	e := &Extension{
		templates: MixinTemplates,
	}

	e.gqlSchemaHooks = []entgql.SchemaHook{addExternalEdges}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

// Templates of the extension
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// GQLSchemaHooks of the extension to seamlessly edit the final gql interface.
func (e *Extension) GQLSchemaHooks() []entgql.SchemaHook {
	return e.gqlSchemaHooks
}

var _ entc.Extension = (*Extension)(nil)
