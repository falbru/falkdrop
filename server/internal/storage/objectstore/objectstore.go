package objectstore

import (
	"context"
)

type ObjectStore interface {
	NewUploadUrl(ctx context.Context, id string) (string, error)
	GetDownloadUrl(ctx context.Context, id string, filename string) (string, error)
	DeleteObject(ctx context.Context, id string) error
	DeleteObjects(ctx context.Context, ids []string) error
}
