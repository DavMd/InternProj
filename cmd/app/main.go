package main

import (
	"InternProj/internal/handler"
	"InternProj/internal/storages"
	"InternProj/internal/storages/memory"
	"InternProj/internal/storages/postgre"
	"log"
	"os"
)

func main() {

	dbType := os.Getenv("DATASTORE_TYPE")

	var store storages.Storage
	var err error

	switch dbType {
	case "memory":
		store = memory.NewMemoryStore()
	case "postgres":
		store, err = postgre.NewPostgreStore()
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
	default:
		log.Fatalf("Unknown database type: %s", dbType)
	}

	handler.ConfigurationHandler(store, "8080")
}
