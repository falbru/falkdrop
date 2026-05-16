package main

import (
	"fmt"
	"os"

	"github.com/falbru/falkdrop/internal/db"
	"github.com/falbru/falkdrop/internal/drop"
)

func main() {
	fmt.Println("Hello World!")

	store := db.NewPostgresStore()
	defer store.Close()
	if store == nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database\n")
	}

	dropRepository := drop.NewPostgresDropRepository(store.Conn)
	dropService := drop.NewDropService(dropRepository)

	err := dropService.CreateDrop([]drop.ResourceType{drop.FileResource, drop.FileResource})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
	}
}
