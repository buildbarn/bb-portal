package graphql

import (
	"github.com/99designs/gqlgen/graphql"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// The Resolver Type for DI
type Resolver struct {
	client *ent.Client
}

// NewSchema creates a graphql executable schema.
func NewSchema(client *ent.Client) graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: &Resolver{client: client},
	})
}
