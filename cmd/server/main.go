package main

import (
	"log"

	"github.com/ARKTEEK/shorty/internal/config"
	"github.com/ARKTEEK/shorty/internal/router"
	"github.com/ARKTEEK/shorty/internal/store"
)

func main() {
	cfg := config.Load()

	database, err := store.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := store.RunMigrations(database, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	r := router.New(database, cfg)

	log.Printf("server listening on %s", cfg.Addr)
	if err := r.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
