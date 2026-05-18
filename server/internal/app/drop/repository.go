package drop

import (
	"time"
)

type DropRepository interface {
	CreateDrop(id DropId, expirationDate time.Time, resourceIds []ResourceId) error
	CreateResource(id ResourceId, resourceType ResourceType) error
	IsUniqueDropId(id DropId) (bool, error)
}
