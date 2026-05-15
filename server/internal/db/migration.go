package db

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

const MIGRATIONS_DIRECTORY = "migrations/"

type Migration struct {
	filePath string
	version  int
}

type SortedMigrations []Migration

func (m SortedMigrations) Len() int {
	return len(m)
}
func (m SortedMigrations) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m SortedMigrations) Less(i, j int) bool {
	return m[i].version < m[j].version
}

func InitMigrations(store *PostgresStore) error {
	rows, err := store.conn.Query(context.Background(), `
		SELECT EXISTS (
		    SELECT 1
		    FROM information_schema.tables
		    WHERE table_schema = 'public'
	        AND table_name = 'schema_migrations'
		);
	`)

	if err != nil {
		return err
	}

	migrationTableExists, err := pgx.CollectOneRow(rows, pgx.RowTo[bool])
	if err != nil {
		return err
	}

	if !migrationTableExists {
		_, err = store.conn.Exec(context.Background(), `
			CREATE TABLE schema_migrations (
				version int PRIMARY KEY
			);
			INSERT INTO schema_migrations (version) VALUES (-1);
		`)

		if err != nil {
			return err
		}
	}

	return nil
}

func NeedsMigration(store *PostgresStore) (bool, error) {
	dbMigrationVersion, err := GetMigrationVersion(store)

	if err != nil {
		return false, err
	}

	migrationFiles, err := getMigrationFiles(MIGRATIONS_DIRECTORY)
	if err != nil {
		return false, err
	}

	return dbMigrationVersion < len(migrationFiles)-1, nil
}

func Migrate(store *PostgresStore) error {
	dbMigrationVersion, err := GetMigrationVersion(store)

	if err != nil {
		return err
	}

	migrationFiles, err := getMigrationFiles(MIGRATIONS_DIRECTORY)
	if err != nil {
		return err
	}

	var migrationErr error
	newMigrationVersion := dbMigrationVersion
	for _, migration := range migrationFiles {
		if migration.version <= dbMigrationVersion {
			continue
		}

		migrationErr = runMigrationFile(store, migration.filePath)
		if migrationErr != nil {
			break
		}

		newMigrationVersion = migration.version
	}

	err = setMigrationVersion(store, newMigrationVersion)
	if migrationErr != nil || err != nil {
		return errors.Join(migrationErr, err)
	}

	return nil
}

func GetMigrationVersion(store *PostgresStore) (int, error) {
	var dbMigrationVersion int
	err := store.conn.QueryRow(context.Background(), "SELECT version FROM schema_migrations").Scan(&dbMigrationVersion)

	if err != nil {
		return -1, err
	}

	return dbMigrationVersion, nil
}

func getMigrationFiles(migrationDirectory string) (SortedMigrations, error) {
	dirEntries, err := os.ReadDir(migrationDirectory)

	if err != nil {
		return nil, err
	}

	migrationFiles := make(SortedMigrations, 0)

	validMigrationFileNameMatcher := regexp.MustCompile("^[\\d]+")

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			return nil, errors.New("The migration directory can't contain subdirectories")
		}

		fileName := dirEntry.Name()

		if !strings.HasSuffix(fileName, ".sql") {
			return nil, errors.New(fmt.Sprintf("Migration file is not a SQL script: '%s'", fileName))
		}

		migrationVersionStr := validMigrationFileNameMatcher.FindString(fileName)
		if migrationVersionStr == "" {
			return nil, errors.New(fmt.Sprintf("Invalid migration name '%s'", fileName))
		}

		migrationVersion, err := strconv.Atoi(migrationVersionStr)
		if err != nil {
			return nil, err
		}

		migrationFiles = append(migrationFiles, Migration{
			filePath: migrationDirectory + fileName,
			version:  migrationVersion,
		})
	}

	sort.Sort(migrationFiles)

	for i, migration := range migrationFiles {
		if i != migration.version {
			return nil, errors.New(fmt.Sprintf("Migration version mismatch '%d' and '%d'", i, migration.version))
		}
	}

	return migrationFiles, nil

}

func runMigrationFile(store *PostgresStore, filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	tx, err := store.conn.Begin(context.Background())
	if err != nil {
		return err
	}

	_, err = tx.Exec(context.Background(), string(contents))
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	return err
}

func setMigrationVersion(store *PostgresStore, version int) error {
	_, err := store.conn.Exec(context.Background(), "UPDATE schema_migrations SET version = $1;", version)
	return err
}
