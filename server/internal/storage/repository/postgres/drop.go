package postgres

import (
	"context"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
)

func (repository PostgresRepository) IsUniqueDropId(id drop.DropId) (bool, error) {
	var dropWithIdExists bool

	err := repository.conn.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM drops WHERE id = $1)", id).Scan(dropWithIdExists)
	if err != nil {
		return false, err
	}

	return !dropWithIdExists, err
}

func (repository PostgresRepository) CreateResource(id drop.ResourceId, resourceType drop.ResourceType) error {
	_, err := repository.conn.Exec(context.Background(), "INSERT INTO resources (id, type) VALUES ($1, $2);", id, resourceType)
	return err
}

func (repository PostgresRepository) CreateDrop(id drop.DropId, expirationDate time.Time, resourceIds []drop.ResourceId) error {
	tx, err := repository.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "INSERT INTO drops (id, expiration_date) VALUES ($1, $2);", id, expirationDate)
	if err != nil {
		return err
	}

	for _, resourceId := range resourceIds {
		_, err = tx.Exec(context.Background(), "INSERT INTO drops_resources (drop_id, resource_id) VALUES ($1, $2);", id, resourceId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	return err
}
