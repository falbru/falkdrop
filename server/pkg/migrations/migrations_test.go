package migrations

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func getMigrationsDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	migrationsDir := filepath.Join(filepath.Dir(filename), "../..", "migrations")
	return filepath.Abs(migrationsDir)
}

func setupTestDatabase(ctx context.Context, t *testing.T) *pgx.Conn {
	ctr, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	testcontainers.CleanupContainer(t, ctr)
	if err != nil {
		t.Fatalf("Failed starting postgres testcontainer: %v", err.Error())
	}

	connStr, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err.Error())
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		t.Fatalf("Failed to connect to postgres database: %v", err.Error())
	}

	return conn
}

func TestMigration(t *testing.T) {
	ctx := context.Background()

	conn := setupTestDatabase(ctx, t)
	defer conn.Close(ctx)

	migrationsDir, err := getMigrationsDir()
	if err != nil {
		t.Fatalf("Can't resolve migrations directory: %v", err.Error())
	}

	t.Run("InitMigrations", func(t *testing.T) {
		err = InitMigrations(conn)
		if err != nil {
			t.Fatalf("Failed to initialize migrations: %v", err.Error())
		}

		migrationVersion, err := GetMigrationVersion(conn)
		if err != nil {
			t.Fatalf("Failed to get migration version: %v", err.Error())
		}

		if migrationVersion != -1 {
			t.Fatalf("Expected migration version to be -1 after initialization, got %d", migrationVersion)
		}
	})

	t.Run("NeedsMigration", func(t *testing.T) {
		migrationFiles, err := getMigrationFiles(migrationsDir)
		if err != nil {
			t.Fatalf("Failed to resolve migration files: %v", err.Error())
		}

		needsMigration, err := NeedsMigration(conn, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to check if migration is needed: %v", err)
		}

		if needsMigration != (len(migrationFiles) > 0) {
			if len(migrationFiles) > 0 {
				t.Fatalf("Expected database to need migration when migration files exist")
			} else {
				t.Fatalf("Expected database not to need migration when no migration files exist")
			}
		}
	})

	t.Run("Migrate", func(t *testing.T) {
		err = Migrate(conn, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to migrate: %v", err.Error())
		}

		needsMigration, err := NeedsMigration(conn, migrationsDir)
		if err != nil {
			t.Fatalf("Failed to check migration is needed after migration: %v", err.Error())
		}
		if needsMigration {
			t.Errorf("Expected database to not need migration after running migrations")
		}
	})
}
