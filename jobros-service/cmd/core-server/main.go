package main

import (
	"github.com/maxime-joseph/Jobros/jobros-service/internal/app"
	"log"
)

func main() {
	_, err := app.NewAppContext()
	if err != nil {
		log.Fatalf("Failed to initialize application context: %v", err)
	}

	// Now you can use appCtx.MongoClient, appCtx.Database, etc.
	log.Println("Application context initialized successfully")
}
