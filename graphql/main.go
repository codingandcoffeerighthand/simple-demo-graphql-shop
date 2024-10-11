package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	AccountUrl string `envconfig:"ACCOUNT_SERVICE_URL"`
	CatalogUrl string `envconfig:"CATALOG_SERVICE_URL"`
	OrderUrl   string `envconfig:"ORDER_SERVICE_URL"`
}

// main runs the GraphQL server.
//
// It reads the ACCOUNT_SERVICE_URL, CATALOG_SERVICE_URL and ORDER_SERVICE_URL
// environment variables, and uses them to construct a GraphQL server that
// delegates to the account, catalog and order services.
//
// It then starts an HTTP server that listens on port 8080 and responds to
// requests for /graphql and /playground.
//
// If the server fails to start (for example, if the environment variables are
// not set), it logs a fatal error and exits.
func main() {
	var cfg AppConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	server, err := NewGraphQLServer(cfg.AccountUrl, cfg.CatalogUrl, cfg.OrderUrl)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/graphql", handler.NewDefaultServer(server.ToExecutableSchema()))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
