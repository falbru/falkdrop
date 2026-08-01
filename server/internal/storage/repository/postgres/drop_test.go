package postgres

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/falbru/falkdrop/pkg/migrations"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func getMigrationsDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", os.ErrNotExist
	}

	migrationsDir := filepath.Join(filepath.Dir(filename), "../../../../", "migrations")
	return filepath.Abs(migrationsDir)
}

func loadFixtures(ctx context.Context, conn *pgx.Conn, t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Failed to get caller info")
	}

	fixturesPath := filepath.Join(filepath.Dir(filename), "testdata", "fixtures.sql")
	fixturesSQL, err := os.ReadFile(fixturesPath)
	if err != nil {
		t.Fatalf("Failed to read fixtures.sql: %v", err.Error())
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err.Error())
	}

	_, err = tx.Exec(ctx, string(fixturesSQL))
	if err != nil {
		t.Fatalf("Failed to execute fixtures: %v", err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.Fatalf("Failed to commit fixtures: %v", err.Error())
	}
}

func setupTestDatabaseWithMigrations(ctx context.Context, t *testing.T) *pgx.Conn {
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

	migrationsDir, err := getMigrationsDir()
	if err != nil {
		t.Fatalf("Can't resolve migrations directory: %v", err.Error())
	}

	err = migrations.InitMigrations(conn)
	if err != nil {
		t.Fatalf("Failed to initialize migrations: %v", err.Error())
	}

	err = migrations.Migrate(conn, migrationsDir)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err.Error())
	}

	loadFixtures(ctx, conn, t)

	return conn
}

func TestPostgresRepository(t *testing.T) {
	ctx := context.Background()

	conn := setupTestDatabaseWithMigrations(ctx, t)
	defer conn.Close(ctx)

	repo := NewPostgresRepositoryFromConnection(conn)

	t.Run("IsUniqueDropId", func(t *testing.T) {
		t.Run("returns true for non-existent drop id", func(t *testing.T) {
			ctx := context.Background()
			dropId := drop.DropId("aaaaa")
			isUnique, err := repo.IsUniqueDropId(ctx, dropId)
			if err != nil {
				t.Fatalf("Failed to check uniqueness: %v", err)
			}
			if !isUnique {
				t.Errorf("Expected drop id (%v) to be unique, got false", dropId)
			}
		})

		t.Run("returns false for existing drop id", func(t *testing.T) {
			ctx := context.Background()
			dropId := drop.DropId("test1")

			isUnique, err := repo.IsUniqueDropId(ctx, dropId)
			if err != nil {
				t.Fatalf("Failed to check uniqueness: %v", err)
			}
			if isUnique {
				t.Errorf("Expected drop id (%v) to not be unique, got true", dropId)
			}
		})
	})

	t.Run("CreateResource", func(t *testing.T) {
		ctx := context.Background()
		err := repo.CreateResource(ctx, drop.ResourceId(uuid.New()), drop.FileResource, nil)

		if err != nil {
			t.Fatalf("Failed to create resource: %v", err.Error())
		}
	})

	t.Run("CreateDrop", func(t *testing.T) {
		ctx := context.Background()
		err := repo.CreateDrop(ctx, drop.DropId("bbbbb"), time.Now(), []drop.ResourceId{})

		if err != nil {
			t.Fatalf("Failed to create drop: %v", err.Error())
		}
	})

	t.Run("GetDropById", func(t *testing.T) {
		ctx := context.Background()
		dropId := drop.DropId("test1")

		drop, err := repo.GetDropById(ctx, dropId)

		if err != nil {
			t.Errorf("Unexpected error when getting drop with drop id %v: %v", dropId, err.Error())
		}

		if drop == nil {
			t.Errorf("Expected to get drop with dropId %v, but result was empty", dropId)
		}

		if year, month, day := drop.ExpirationDate.Date(); year != 2050 || month != time.January || day != 1 {
			t.Errorf("Expected date of drop %v to be 2050-01-01, but got %v", dropId, drop.ExpirationDate.String())
		}

		if len(drop.Resources) != 2 {
			t.Errorf("Expected resources length of drop %v to be 2, but got %v", dropId, len(drop.Resources))
		}

		expectedResourceNames := []string{"file1.txt", "file2.txt"}

		resourceNames := make([]string, len(drop.Resources))
		for i, resource := range drop.Resources {
			if resource.Name == nil {
				resourceNames[i] = ""
			} else {
				resourceNames[i] = *resource.Name
			}
		}

		sort.Strings(expectedResourceNames)
		sort.Strings(resourceNames)

		for i, resourceName := range resourceNames {
			if expectedResourceNames[i] != resourceName {
				t.Errorf("Expected resource name to be %v, but got %v", expectedResourceNames[i], resourceName)
			}
		}
	})

	t.Run("GetResourcesByDropId", func(t *testing.T) {
		ctx := context.Background()
		dropId := drop.DropId("test1")

		resources, err := repo.GetResourcesByDropId(ctx, dropId)
		if err != nil {
			t.Fatalf("Failed to get resources by dropId %v: %v", dropId, err.Error())
		}

		if len(resources) != 2 {
			t.Fatalf("Expected resources length of drop %v to be 2, but got %v", dropId, len(resources))
		}
		idSet := make(map[drop.ResourceId]bool, 2)
		idSet[drop.ResourceId(uuid.MustParse("11111111-1111-1111-1111-111111111111"))] = true
		idSet[drop.ResourceId(uuid.MustParse("22222222-2222-2222-2222-222222222222"))] = true

		for _, resource := range resources {
			if !idSet[resource.Id] {
				t.Errorf("Unexpected resource ID: %v", resource.Id)
			}
		}
	})

	t.Run("GetDropsExpiredByDate", func(t *testing.T) {
		ctx := context.Background()

		drops, err := repo.GetDropsExpiredByDate(ctx, time.Date(2026, time.January, 1, 1, 1, 1, 1, time.UTC))
		if err != nil {
			t.Fatalf("Failed to get expired drops after date: %v", err.Error())
		}

		if len(drops) != 2 {
			t.Fatalf("Expected to get two expired drops, but got: %v", len(drops))
		}

		for _, drop := range drops {
			if drop.Id == "expd1" {
				if drop.ExpirationDate.Year() != 2000 || drop.ExpirationDate.Month() != 1 || drop.ExpirationDate.Day() != 1 {
					t.Fatalf("Expected expirationDate for expd1 to be 2000-01-01, but got: %v", drop.ExpirationDate.Format("2000-01-01"))
				}

				if len(drop.Resources) != 2 {
					t.Fatalf("Expected expd1 to have 2 resources, but got: %v", len(drop.Resources))
				}
			} else if drop.Id == "expd2" {
				if drop.ExpirationDate.Year() != 2005 || drop.ExpirationDate.Month() != 1 || drop.ExpirationDate.Day() != 1 {
					t.Fatalf("Expected expirationDate for expd2 to be 2000-01-01, but got: %v", drop.ExpirationDate.Format("2000-01-01"))
				}

				if len(drop.Resources) != 1 {
					t.Fatalf("Expected expd2 to have 2 resources, but got: %v", len(drop.Resources))
				}
			} else {
				t.Fatalf("Did not expect to get drop: %v", drop.Id)
			}
		}
	})

	t.Run("GetResourcesByIds", func(t *testing.T) {
		ctx := context.Background()
		resourceIds := []drop.ResourceId{drop.ResourceId(uuid.MustParse("11111111-1111-1111-1111-111111111111")), drop.ResourceId(uuid.MustParse("22222222-2222-2222-2222-222222222222"))}

		resources, err := repo.GetResourcesByIds(ctx, resourceIds)
		if err != nil {
			t.Fatalf("Failed to get resources by ids: %v", err.Error())
		}

		if len(resources) != len(resourceIds) {
			t.Fatalf("Expected resources length to be %v, but got %v", len(resourceIds), len(resources))
		}

		idSet := make(map[drop.ResourceId]bool, len(resourceIds))
		for _, id := range resourceIds {
			idSet[id] = true
		}

		for _, resource := range resources {
			if !idSet[resource.Id] {
				t.Errorf("Unexpected resource ID: %v", resource.Id)
			}
		}
	})
}
