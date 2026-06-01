package mock

import (
	"context"
)

type MockObjectStore struct {
	NewUploadUrlFn   func(ctx context.Context, id string) (string, error)
	GetDownloadUrlFn func(ctx context.Context, id string, filename string) (string, error)
}

func (objectStore MockObjectStore) NewUploadUrl(ctx context.Context, id string) (string, error) {
	if objectStore.NewUploadUrlFn != nil {
		return objectStore.NewUploadUrlFn(ctx, id)
	}

	return "http://example.org", nil
}

func (objectStore MockObjectStore) GetDownloadUrl(ctx context.Context, id string, filename string) (string, error) {
	if objectStore.GetDownloadUrlFn != nil {
		return objectStore.GetDownloadUrlFn(ctx, id, filename)
	}

	return "http://example.org", nil
}

func NewMockObjectStore() MockObjectStore {
	return MockObjectStore{}
}

func (objectstore MockObjectStore) WithNewUploadUrl(fn func(ctx context.Context, id string) (string, error)) MockObjectStore {
	objectstore.NewUploadUrlFn = fn
	return objectstore
}

func (objectstore MockObjectStore) WithGetDownloadUrl(fn func(ctx context.Context, id string, filename string) (string, error)) MockObjectStore {
	objectstore.GetDownloadUrlFn = fn
	return objectstore
}
