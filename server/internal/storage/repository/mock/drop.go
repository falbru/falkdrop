package mock

import (
	"context"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
)

type MockDropRepository struct {
	getResourcesByIdsFn  func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error)
	getDropByIdFn        func(ctx context.Context, id drop.DropId) (*drop.Drop, error)
	withCreateResourceFn func(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error
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

func (repo MockDropRepository) WithGetResourcesByIds(fn func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error)) MockDropRepository {
	repo.getResourcesByIdsFn = fn
	return repo
}

func (repo MockDropRepository) WithGetDropById(fn func(ctx context.Context, id drop.DropId) (*drop.Drop, error)) MockDropRepository {
	repo.getDropByIdFn = fn
	return repo
}

func (repo MockDropRepository) WithCreateResource(fn func(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error) MockDropRepository {
	repo.withCreateResourceFn = fn
	return repo
}
