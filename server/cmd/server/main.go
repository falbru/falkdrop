package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/handlers"
	"github.com/falbru/falkdrop/internal/storage/objectstore/s3"
	"github.com/falbru/falkdrop/internal/storage/repository/postgres"
)

func main() {
	repository, err := postgres.NewPostgresRepository(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database: %s\n", err.Error())
		os.Exit(1)
	}
	defer repository.Close(context.Background())

	objectStore, err := s3.NewS3Store(context.Background(), "falkdrop2")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to object store: %s\n", err.Error())
		os.Exit(1)
	}

	dropService := drop.NewDropService(repository, objectStore)

	resourceHandler := handlers.NewResourceHandler(dropService)
	dropHandler := handlers.NewDropHandler(dropService)

	http.HandleFunc("POST /resource", resourceHandler.Create)
	http.HandleFunc("POST /drop", dropHandler.Create)
	http.HandleFunc("GET /drop/{dropId}", dropHandler.Get)

	http.ListenAndServe(":8080", nil)
}
