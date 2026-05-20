package main

import (
	"context"
	"fmt"
	"os"

	"github.com/falbru/falkdrop/internal/app/drop"
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

	resourceId, url, err := dropService.CreateResourceWithUploadLink(drop.FileResource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: could not create presigned url: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(url)

	_, err = dropService.CreateDrop([]drop.ResourceId{resourceId})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: create drop: %v\n", err.Error())
		os.Exit(1)
	}

	dropWithResourceLinks, err := dropService.GetDropWithResourceLinks("xf0lx")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: get drop with resource links: %v\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(dropWithResourceLinks.Id)
	fmt.Println(dropWithResourceLinks.ExpirationDate)
	for _, resource := range dropWithResourceLinks.ResourceLinks {
		fmt.Println(resource.Link)
	}
}
