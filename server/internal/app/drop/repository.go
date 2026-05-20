package drop

import (
	"context"
	"time"
)

type DropRepository interface {
	CreateDrop(ctx context.Context, id DropId, expirationDate time.Time, resourceIds []ResourceId) error
	CreateResource(ctx context.Context, id ResourceId, resourceType ResourceType) error
	GetDropById(ctx context.Context, id DropId) (*Drop, error)
	GetResourcesByDropId(ctx context.Context, id DropId) ([]Resource, error)
	IsUniqueDropId(ctx context.Context, id DropId) (bool, error)
}
