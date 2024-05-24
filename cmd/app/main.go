package main

import (
	"InternProj/internal/handler"
	"InternProj/internal/storages"
	"InternProj/internal/storages/memory"
	"InternProj/internal/storages/postgre"
	"flag"
	"log"
)

func main() {

	var dbType string
	var connString string

	flag.StringVar(&dbType, "db", "memory", "Database type: 'memory' or 'postgres'")
	flag.StringVar(&connString, "conn", "", "Connection string for PostgreSQL")
	flag.Parse()

	var store storages.Storage
	var err error

	switch dbType {
	case "memory":
		store = memory.NewMemoryStore()
	case "postgres":
		if connString == "" {
			log.Fatal("Connection string is required for PostgreSQL")
		}
		store, err = postgre.NewPostgreStore(connString)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
	default:
		log.Fatalf("Unknown database type: %s", dbType)
	}

	handler.ConfigurationHandler(store, "8080")
}
