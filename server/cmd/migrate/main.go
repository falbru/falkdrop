package main

import (
	"fmt"
	"os"

	"github.com/falbru/falkdrop/internal/db"
	"github.com/falbru/falkdrop/pkg/migrations"
)

func printHelp() {
	fmt.Println(`Usage: migrate <command>

Commands:
  check    Check if pending migrations exist (exit code 1 if yes, 0 if up-to-date, 2 if error occurred)
  up       Apply all pending migrations`)
}

const MIGRATIONS_DIRECTORY = "migrations/"

func main() {
	args := os.Args[1:]

	if len(args) == 0 || (args[0] != "check" && args[0] != "up") {
		printHelp()
		os.Exit(1)
	}

	store := db.NewPostgresStore()
	if store == nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database")
		os.Exit(2)
	}

	migrations.InitMigrations(store.Conn)

	switch args[0] {
	case "check":
		needsMigration, err := migrations.NeedsMigration(store.Conn, MIGRATIONS_DIRECTORY)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to check migrations: %v\n", err)
			os.Exit(2)
		}
		if needsMigration {
			os.Exit(1)
		}
		os.Exit(0)
	case "up":
		needsMigration, err := migrations.NeedsMigration(store.Conn, MIGRATIONS_DIRECTORY)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to check migrations: %v\n", err)
			os.Exit(2)
		}

		if needsMigration {
			initialVersion, err := migrations.GetMigrationVersion(store.Conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to check initial migration version: %v\n", err)
				os.Exit(2)
			}

			err = migrations.Migrate(store.Conn, MIGRATIONS_DIRECTORY)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to apply migrations: %v\n", err)
				os.Exit(2)
			}

			finalVersion, err := migrations.GetMigrationVersion(store.Conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to check final migration version: %v\n", err)
				os.Exit(2)
			}

			fmt.Printf("Migrated from version %d to %d\n", initialVersion, finalVersion)
			fmt.Println("Migrated successfully!")
		} else {
			fmt.Println("Already up-to-date")
		}

	default:
		printHelp()
		os.Exit(1)
	}
}
