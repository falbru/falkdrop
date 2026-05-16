package drop

import "time"

type ResourceType string

const (
	FileResource ResourceType = "file"
	TextResource              = "text"
)

type Resource struct {
	Type ResourceType
	Link string
}

type DropId string

type Drop struct {
	Id             DropId
	ExpirationDate time.Time
	Resources      []Resource
}
