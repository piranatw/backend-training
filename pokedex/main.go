package main

import (
	"log"
	"net/http"
	"os"
	"pokedex/database"
	"pokedex/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const defaultPort = "8070"

func main() {
	port := getEnv("PORT", defaultPort)
	dbPath := getEnv("POKEDEX_DB_PATH", "pokedex.db?_foreign_keys=on")

	db, err := database.Connect(dbPath)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	pokemonRepo := database.NewPokemonRepository(db)
	resolver := &graph.Resolver{
		PokemonRepo: pokemonRepo,
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	server := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: resolver,
	}))

	router.Handle("/", playground.Handler("Pokedex GraphQL", "/query"))
	router.Handle("/query", server)

	log.Printf("GraphQL playground: http://localhost:%s/", port)
	log.Printf("GraphQL endpoint: http://localhost:%s/query", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
