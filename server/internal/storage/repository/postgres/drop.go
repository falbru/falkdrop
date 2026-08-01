package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (repository PostgresRepository) IsUniqueDropId(ctx context.Context, id drop.DropId) (bool, error) {
	var dropWithIdExists bool

	err := repository.conn.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM drops WHERE id = $1)", id).Scan(&dropWithIdExists)
	if err != nil {
		return false, err
	}

	return !dropWithIdExists, err
}

func (repository PostgresRepository) CreateResource(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error {
	_, err := repository.conn.Exec(ctx, "INSERT INTO resources (id, type, name) VALUES ($1, $2, $3);", id, resourceType, name)
	return err
}

func (repository PostgresRepository) CreateDrop(ctx context.Context, id drop.DropId, expirationDate time.Time, resourceIds []drop.ResourceId) error {
	tx, err := repository.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "INSERT INTO drops (id, expiration_date) VALUES ($1, $2);", id, expirationDate)
	if err != nil {
		return err
	}

	for _, resourceId := range resourceIds {
		_, err = tx.Exec(ctx, "UPDATE resources SET drop_id = $1 WHERE id = $2;", id, resourceId)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	return err
}

func (repository PostgresRepository) GetDropById(ctx context.Context, id drop.DropId) (*drop.Drop, error) {
	var dropId string
	var expirationDate time.Time
	err := repository.conn.QueryRow(ctx, "SELECT id, expiration_date FROM drops WHERE id = $1", id).Scan(&dropId, &expirationDate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	resources, err := repository.GetResourcesByDropId(ctx, id)

	if err != nil {
		return nil, err
	}

	return &drop.Drop{
		Id:             drop.DropId(dropId),
		ExpirationDate: expirationDate,
		Resources:      resources,
	}, err
}

func (repository PostgresRepository) GetDropsExpiredByDate(ctx context.Context, date time.Time) ([]drop.Drop, error) {
	queryRows, err := repository.conn.Query(ctx, "SELECT drops.id, drops.expiration_date, resources.id, resources.type, resources.name FROM drops LEFT JOIN resources ON resources.drop_id = drops.id WHERE drops.expiration_date <= $1 ORDER BY drops.id", date)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []drop.Drop{}, nil
		}

		return nil, err
	}

	var expiredDrops []drop.Drop

	var previousDropId drop.DropId

	var dropId drop.DropId
	var expirationDate time.Time
	var resourceId *string
	var resourceType *string
	var resourceName *string
	_, err = pgx.ForEachRow(queryRows, []any{&dropId, &expirationDate, &resourceId, &resourceType, &resourceName}, func() error {
		if len(expiredDrops) == 0 || previousDropId != dropId {
			previousDropId = dropId
			expiredDrops = append(expiredDrops, drop.Drop{
				Id:             dropId,
				ExpirationDate: expirationDate,
				Resources:      []drop.Resource{},
			})
		}

		if resourceId != nil && resourceType != nil {
			resource := drop.Resource{
				Id:     drop.ResourceId(uuid.MustParse(*resourceId)), // TODO must parse is unsafe
				Type:   drop.ResourceType(*resourceType),
				Name:   resourceName,
				DropId: &dropId,
			}

			expiredDrops[len(expiredDrops)-1].Resources = append(expiredDrops[len(expiredDrops)-1].Resources, resource)
		}

		return nil
	})

	return expiredDrops, err
}

func (repository PostgresRepository) DeleteDropsById(ctx context.Context, ids []drop.DropId) error {
	idStrs := make([]string, len(ids))

	for i, id := range ids {
		idStrs[i] = string(id)
	}

	_, err := repository.conn.Exec(ctx, "DELETE FROM drops WHERE id = ANY($1)", idStrs)
	return err
}

func (repository PostgresRepository) GetResourcesByDropId(ctx context.Context, id drop.DropId) ([]drop.Resource, error) {
	queryRows, err := repository.conn.Query(ctx, "SELECT id, type, name FROM resources WHERE drop_id = $1", id)
	if err != nil {
		return []drop.Resource{}, err
	}

	rows, err := pgx.CollectRows(queryRows, func(row pgx.CollectableRow) (drop.Resource, error) {
		var resourceId uuid.UUID
		var resourceType drop.ResourceType
		var resourceName *string
		err := row.Scan(&resourceId, &resourceType, &resourceName)

		return drop.Resource{
			Id:   drop.ResourceId(resourceId),
			Type: resourceType,
			Name: resourceName,
		}, err
	})

	return rows, err
}

func (repository PostgresRepository) GetResourcesByIds(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error) {
	if len(ids) == 0 {
		return []drop.Resource{}, nil
	}

	idUUIDs := make([]uuid.UUID, len(ids))
	for i, id := range ids {
		idUUIDs[i] = uuid.UUID(id)
	}

	queryRows, err := repository.conn.Query(ctx, "SELECT id, type, name, drop_id FROM resources WHERE id = ANY($1)", idUUIDs)
	if err != nil {
		return []drop.Resource{}, err
	}

	rows, err := pgx.CollectRows(queryRows, func(row pgx.CollectableRow) (drop.Resource, error) {
		var resourceId uuid.UUID
		var resourceType drop.ResourceType
		var resourceName *string
		var resourceDropId *drop.DropId
		err := row.Scan(&resourceId, &resourceType, &resourceName, &resourceDropId)

		return drop.Resource{
			Id:     drop.ResourceId(resourceId),
			Type:   resourceType,
			Name:   resourceName,
			DropId: resourceDropId,
		}, err
	})

	return rows, err
}
