package main

import (
	"context"
	"log"
	"net/http"

	"social-network/db/sqlite"
	"social-network/internal/config"
	"social-network/internal/handlers"
	"social-network/internal/repository"
	"social-network/internal/service"
	chatsvc "social-network/internal/service/chat"
	"social-network/internal/utils"
)

//--------------------------------------------------------------------------------------|

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

	// Initialize WebSocket Hub
	hub := chatsvc.NewHub(rep.Chat, rep.User, rep.Group)
	go hub.Run(context.Background())

	serv := service.NewService(rep, hub)

	// Initialize Uploader
	uploadDir := "./uploads"
	uploadURL := "/api/v1/uploads"
	uploader := utils.NewLocalImageUploader(uploadDir, uploadURL)

	hdr := handlers.NewHandler(serv, hub, rep.User, rep.Chat, uploader)
	mux := config.SetupRoutes(hdr)

	// Serve uploaded images
	mux.Handle(uploadURL+"/", http.StripPrefix(uploadURL, http.FileServer(http.Dir(uploadDir))))

	port := ":8080"
	log.Println("Starting server on ", port)
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
