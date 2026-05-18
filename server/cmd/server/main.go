package main

import (
	"fmt"
	"os"

	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/internal/storage/objectstore/s3"
	"github.com/falbru/falkdrop/internal/storage/repository/postgres"
)

func main() {
	fmt.Println("Hello World!")

	repository, err := postgres.NewPostgresRepository()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database\n")
		os.Exit(1)
	}
	defer repository.Close()

	objectStore := s3.NewGarageStore("drops")

	dropService := drop.NewDropService(repository, objectStore)

	resourceId, url, err := dropService.CreateResourceWithUploadLink(drop.FileResource)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: could not create presigned url: ", err.Error())
	}

	fmt.Println(url)

	err = dropService.CreateDrop([]drop.ResourceId{resourceId})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
	}
}
