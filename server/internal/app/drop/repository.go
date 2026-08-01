package drop

import (
	"context"
	"time"
)

type DropRepository interface {
	CreateDrop(ctx context.Context, id DropId, expirationDate time.Time, resourceIds []ResourceId) error
	CreateResource(ctx context.Context, id ResourceId, resourceType ResourceType, name *string) error
	GetDropById(ctx context.Context, id DropId) (*Drop, error)
	GetDropsExpiredByDate(ctx context.Context, date time.Time) ([]Drop, error)
	DeleteDropsById(ctx context.Context, id []DropId) error
	GetResourcesByDropId(ctx context.Context, id DropId) ([]Resource, error)
	GetResourcesByIds(ctx context.Context, ids []ResourceId) ([]Resource, error)
	IsUniqueDropId(ctx context.Context, id DropId) (bool, error)
}
