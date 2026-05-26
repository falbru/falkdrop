package objectstore

import (
	"context"
)

type ObjectStore interface {
	NewUploadUrl(ctx context.Context, id string) (string, error)
	GetDownloadUrl(ctx context.Context, id string, filename string) (string, error)
}
