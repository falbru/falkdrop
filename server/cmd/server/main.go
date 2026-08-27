package main

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/falbru/falkdrop/internal/api/handlers"
	"github.com/falbru/falkdrop/internal/api/middleware"
	authService "github.com/falbru/falkdrop/internal/app/auth"
	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/auth"
	"github.com/falbru/falkdrop/internal/storage/objectstore"
	"github.com/falbru/falkdrop/internal/storage/objectstore/s3"
	"github.com/falbru/falkdrop/internal/storage/repository/postgres"
	"github.com/go-co-op/gocron/v2"
	"github.com/joho/godotenv"

	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not load environment variables from .env file: %s\n", err.Error())
		os.Exit(1)
	}

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

	keycloakRealmUrl, isKeycloakRealmUrlSet := os.LookupEnv("KEYCLOAK_REALM_URL")
	if !isKeycloakRealmUrlSet {
		fmt.Fprintf(os.Stderr, "Error: KEYCLOAK_REALM_URL not set\n")
	}

	keycloakClientId, isKeycloakClientSet := os.LookupEnv("KEYCLOAK_CLIENT_ID")
	if !isKeycloakClientSet {
		fmt.Fprintf(os.Stderr, "Error: KEYCLOAK_CLIENT_ID not set\n")
	}

	provider := auth.NewAuthProvider(keycloakRealmUrl, keycloakClientId)
	err = provider.Init(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create auth provider: %s", err.Error())
		os.Exit(1)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to setup cron scheduler: %s", err.Error())
		os.Exit(1)
	}

	authService := authService.NewAuthService(provider)
	dropService := drop.NewDropService(repository, objectStore, authService)

	resourceHandler := handlers.NewResourceHandler(dropService)
	dropHandler := handlers.NewDropHandler(dropService)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /resource", middleware.ErrorHandler(resourceHandler.Create))
	mux.HandleFunc("POST /drop", middleware.ErrorHandler(dropHandler.Create))
	mux.HandleFunc("GET /drop/{dropId}", middleware.ErrorHandler(dropHandler.Get))

	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	c := cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}

	job, err := drop.RegisterJobs(scheduler, dropService)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to setup drop jobs: %s", err.Error())
		os.Exit(1)
	}
	slog.Info("started drop job", "ID", job.ID().String())

	scheduler.Start()
	defer scheduler.Shutdown()

	slog.Info("Server starting on http://localhost:8082")
	http.Serve(l, c.Handler(mux))
}
