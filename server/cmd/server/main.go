package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/falbru/falkdrop/internal/api/handlers"
	"github.com/falbru/falkdrop/internal/api/middleware"
	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/storage/objectstore"
	"github.com/falbru/falkdrop/internal/storage/objectstore/s3"
	"github.com/falbru/falkdrop/internal/storage/repository/postgres"

	"github.com/rs/cors"
)

func main() {
	repository, err := postgres.NewPostgresRepository(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database: %s\n", err.Error())
		os.Exit(1)
	}
	defer repository.Close(context.Background())

	var objectStore objectstore.ObjectStore
	switch os.Getenv("OBJECTSTORE_PROVIDER") {
	case "s3":
		bucketName, isBucketNameSet := os.LookupEnv("S3_BUCKET")
		if !isBucketNameSet {
			fmt.Fprintf(os.Stderr, "Error: S3_BUCKET not set\n")
			os.Exit(1)
		}

		objectStore, err = s3.NewS3Store(context.Background(), bucketName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not connect to object store: %s\n", err.Error())
			os.Exit(1)
		}
	case "garage":
		bucketName, isBucketNameSet := os.LookupEnv("GARAGE_BUCKET")
		if !isBucketNameSet {
			fmt.Fprintf(os.Stderr, "Error: GARAGE_BUCKET not set\n")
			os.Exit(1)
		}

		objectStore, err = s3.NewGarageStore(context.Background(), bucketName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Could not connect to object store: %s\n", err.Error())
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: OBJECTSTORE_PROVIDER set to invalid value")
		os.Exit(1)
	}

	dropService := drop.NewDropService(repository, objectStore)

	resourceHandler := handlers.NewResourceHandler(dropService)
	dropHandler := handlers.NewDropHandler(dropService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /resource", middleware.ErrorHandler(resourceHandler.Create))
	mux.HandleFunc("POST /drop", middleware.ErrorHandler(dropHandler.Create))
	mux.HandleFunc("GET /drop/{dropId}", middleware.ErrorHandler(dropHandler.Get))

	c := cors.AllowAll()

	http.ListenAndServe(":8080", c.Handler(mux))
}
