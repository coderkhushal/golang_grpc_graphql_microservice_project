package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/coderkhushal/go-grpc-graphql-microservices/account"
	"github.com/coderkhushal/go-grpc-graphql-microservices/catalog"
	"github.com/coderkhushal/go-grpc-graphql-microservices/order"
)

type Server struct {
	accoutClient  *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQLServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	accountClient, err := account.NewClient(accountUrl)
	if err != nil {
		return nil, err
	}

	catalogClient, err := catalog.NewClient(catalogUrl)
	if err != nil {
		accountClient.Close()
		return nil, err
	}
	orderClient, err := order.NewClient(orderUrl)

	if err != nil {
		orderClient.Close()
		accountClient.Close()
		return nil, err
	}
	return &Server{
		accoutClient:  accountClient,
		orderClient:   orderClient,
		catalogClient: catalogClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &Mutationresolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryresolver{
		server: s,
	}
}
func (s *Server) Account() AccountResolver {
	return &Accountresolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(Config{
		Resolvers: s,
	})
}
