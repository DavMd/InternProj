package handler

import (
	"InternProj/graph"
	"InternProj/graph/generated"
	"InternProj/internal/storages"
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func ConfigurationHandler(store storages.Storage, port string) {

	resolver := &graph.Resolver{Store: store}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
