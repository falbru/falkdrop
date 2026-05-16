package drop

import "time"

type ResourceType string

const (
	FileResource ResourceType = "file"
	TextResource = "text"
)

type Resource struct {
	Type ResourceType
	Link string
}

type Drop struct {
	Id             string
	ExpirationDate time.Time
	Resources      []Resource
}
