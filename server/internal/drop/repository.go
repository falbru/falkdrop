package drop

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type DropRepository interface {
	CreateDrop(id DropId, expirationDate time.Time, resources []Resource) error
	IsUniqueDropId(id DropId) (bool, error)
}

type PostgresDropRepository struct {
	conn *pgx.Conn
}

func NewPostgresDropRepository(conn *pgx.Conn) *PostgresDropRepository {
	return &PostgresDropRepository{
		conn,
	}
}

func (repository PostgresDropRepository) IsUniqueDropId(id DropId) (bool, error) {
	var dropWithIdExists bool

	err := repository.conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM drops WHERE id = $1)", id).Scan(dropWithIdExists)
	if err != nil {
		return false, err
	}

	return !dropWithIdExists, err
}

func (repository PostgresDropRepository) CreateDrop(id DropId, expirationDate time.Time, resources []Resource) error {
	tx, err := repository.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "INSERT INTO drops (id, expiration_date) VALUES ($1, $2);", id, expirationDate)
	if err != nil {
		return err
	}

	for _, resource := range resources {
		_, err = tx.Exec(context.Background(), "INSERT INTO resources (link, type) VALUES ($1, $2);", resource.Link, resource.Type)
		if err != nil {
			return err
		}

		_, err = tx.Exec(context.Background(), "INSERT INTO drops_resources (drop_id, resource_link) VALUES ($1, $2);", id, resource.Link)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	return err
}
