package mock

import (
	"context"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
)

type MockDropRepository struct {
	getResourcesByIdsFn     func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error)
	getDropByIdFn           func(ctx context.Context, id drop.DropId) (*drop.Drop, error)
	getDropsExpiredByDateFn func(ctx context.Context, date time.Time) ([]drop.Drop, error)
	withCreateResourceFn    func(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error
	deleteDropsByIdFn       func(ctx context.Context, ids []drop.DropId) error
}

func NewMockDropRepository() MockDropRepository {
	return MockDropRepository{}
}

func (repo MockDropRepository) CreateDrop(ctx context.Context, id drop.DropId, expirationDate time.Time, resourceIds []drop.ResourceId) error {
	return nil
}

func (repo MockDropRepository) CreateResource(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error {
	return nil
}

func (repo MockDropRepository) GetDropById(ctx context.Context, id drop.DropId) (*drop.Drop, error) {
	if repo.getDropByIdFn != nil {
		return repo.getDropByIdFn(ctx, id)
	}

	return nil, nil
}

func (repo MockDropRepository) GetDropsExpiredByDate(ctx context.Context, date time.Time) ([]drop.Drop, error) {
	if repo.getDropsExpiredByDateFn != nil {
		return repo.getDropsExpiredByDateFn(ctx, date)
	}

	return []drop.Drop{}, nil
}

func (repo MockDropRepository) DeleteDropsById(ctx context.Context, ids []drop.DropId) error {
	if repo.deleteDropsByIdFn != nil {
		return repo.deleteDropsByIdFn(ctx, ids)
	}

	return nil
}

func (repo MockDropRepository) GetResourcesByDropId(ctx context.Context, id drop.DropId) ([]drop.Resource, error) {
	return []drop.Resource{}, nil
}

func (repo MockDropRepository) GetResourcesByIds(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error) {
	if repo.getResourcesByIdsFn != nil {
		return repo.getResourcesByIdsFn(ctx, ids)
	}

	return []drop.Resource{}, nil
}

func (repo MockDropRepository) IsUniqueDropId(ctx context.Context, id drop.DropId) (bool, error) {
	return true, nil

}

func (repo MockDropRepository) WithGetDropById(fn func(ctx context.Context, id drop.DropId) (*drop.Drop, error)) MockDropRepository {
	repo.getDropByIdFn = fn
	return repo
}

func (repo MockDropRepository) WithGetResourcesByIds(fn func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error)) MockDropRepository {
	repo.getResourcesByIdsFn = fn
	return repo
}

func (repo MockDropRepository) WithGetDropsExpiredByDate(fn func(ctx context.Context, date time.Time) ([]drop.Drop, error)) MockDropRepository {
	repo.getDropsExpiredByDateFn = fn
	return repo
}

func (repo MockDropRepository) WithDeleteDropsById(fn func(ctx context.Context, ids []drop.DropId) error) MockDropRepository {
	repo.deleteDropsByIdFn = fn
	return repo
}

func (repo MockDropRepository) WithCreateResource(fn func(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error) MockDropRepository {
	repo.withCreateResourceFn = fn
	return repo
}
