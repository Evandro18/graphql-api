package gql

import (
	"../postgres"
	"github.com/graphql-go/graphql"
)

type Root struct {
	Query *graphql.Object
}

func NewRoot(db *postgres.Db) *Root  {
	// Create a resolver holding our database. Resolver can be found in resolvers.go
	resolver := Resolver{db: db}

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called name
	root := Root{
		Query: graphql.ObjectConfig(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Field{
					// Slice of User type which can be found in types.go
					Type: graphql.NewList(User)
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type: graphql.String,
						},
					},
					Resolve: resolver.UserResolver,
				},
			},
		),
	}
	return &root
}