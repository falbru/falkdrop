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

func (resourceId ResourceId) String() string {
	return uuid.UUID(resourceId).String()
}

type Resource struct {
	Id   ResourceId
	Type ResourceType
	Name *string
}

type ResourceWithUploadUrl struct {
	Resource
	UploadUrl string
}

type DropId string

type Drop struct {
	Id             DropId
	ExpirationDate time.Time
	Resources      []Resource
}

type ResourceWithDownloadUrl struct {
	Resource
	DownloadUrl string
}

type DropWithResourceDownloadUrls struct {
	Id             DropId
	ExpirationDate time.Time
	Resources      []ResourceWithDownloadUrl
}
