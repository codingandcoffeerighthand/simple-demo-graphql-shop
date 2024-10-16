//go:generate go run github.com/99designs/gqlgen generate
package main

import (
	"fmt"
	"shop-graphql-demo/account"
	"shop-graphql-demo/catalog"
	"shop-graphql-demo/order"

	"github.com/99designs/gqlgen/graphql"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(
	accountUrl, catalogUrl, orderUrl string,
) (*Server, error) {

	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to create account client: %w", err)
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, fmt.Errorf("failed to create catalog client: %w", err)
	}

	orderClient, err := order.NewClient(orderUrl)
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		return nil, fmt.Errorf("failed to create order client: %w", err)
	}
	return &Server{accountClient, catalogClient, orderClient}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{server: s}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{server: s}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{server: s}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{Resolvers: s})
}
