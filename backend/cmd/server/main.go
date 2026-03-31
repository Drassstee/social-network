package main

import (
	"log"
	"net/http"

	"social-network/db/sqlite"
	"social-network/internal/config"
	"social-network/internal/handlers"
	"social-network/internal/repository"
	"social-network/internal/service"
)

func main() {
	dbPath := "social_network.db"
	db, err := sqlite.ConnectDB(dbPath)
	if err != nil {
		log.Fatal("database connection failed: ", err)
	}
	defer db.Close()

	migPath := "db/migrations/sqlite"
	if err := sqlite.RunMigrations(db, migPath); err != nil {
		log.Fatal("Migrations failed: ", err)
	}

	rep := repository.NewRepo(db)
	serv := service.NewService(rep)
	hdr := handlers.NewHandler(serv)
	mux := config.SetupRoutes(hdr)

	port := ":8080"
	log.Println("Starting server on ", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
