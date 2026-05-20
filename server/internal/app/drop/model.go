package drop

import (
	"time"

	"github.com/google/uuid"
)

type ResourceType string

const (
	FileResource ResourceType = "file"
	TextResource              = "text"
)

type ResourceId uuid.UUID

func (resouceId ResourceId) String() string {
	return uuid.UUID(resouceId).String()
}

type Resource struct {
	Id   ResourceId
	Type ResourceType
}

type DropId string

type Drop struct {
	Id             DropId
	ExpirationDate time.Time
	Resources      []Resource
}

type ResourceLink struct {
	Type ResourceType
	Link string
}

type DropWithResourceLinks struct {
	Id             DropId
	ExpirationDate time.Time
	ResourceLinks  []ResourceLink
}
