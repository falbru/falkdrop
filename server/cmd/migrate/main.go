package main

import (
	"context"
	"fmt"
	"os"

	"github.com/falbru/falkdrop/pkg/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func printHelp() {
	fmt.Println(`Usage: migrate <command>

Commands:
  check    Check if pending migrations exist (exit code 1 if yes, 0 if up-to-date, 2 if error occurred)
  up       Apply all pending migrations`)
}

const MIGRATIONS_DIRECTORY = "migrations/"

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not load environment variables from .env file: %s\n", err.Error())
		os.Exit(1)
	}

	args := os.Args[1:]

	if len(args) == 0 || (args[0] != "check" && args[0] != "up") {
		printHelp()
		os.Exit(1)
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not connect to database")
		os.Exit(2)
	}

	err = migrations.InitMigrations(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to initialize migrations")
		os.Exit(2)
	}

	switch args[0] {
	case "check":
		needsMigration, err := migrations.NeedsMigration(conn, MIGRATIONS_DIRECTORY)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to check migrations: %v\n", err)
			os.Exit(2)
		}
		if needsMigration {
			os.Exit(1)
		}
		os.Exit(0)
	case "up":
		needsMigration, err := migrations.NeedsMigration(conn, MIGRATIONS_DIRECTORY)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to check migrations: %v\n", err)
			os.Exit(2)
		}

		if needsMigration {
			initialVersion, err := migrations.GetMigrationVersion(conn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to check initial migration version: %v\n", err)
				os.Exit(2)
			}

			err = migrations.Migrate(conn, MIGRATIONS_DIRECTORY)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: failed to apply migrations: %v\n", err)
				os.Exit(2)
			}

			finalVersion, err := migrations.GetMigrationVersion(conn)
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
