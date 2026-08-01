package mock

import (
	"context"
)

type MockObjectStore struct {
	newUploadUrlFn   func(ctx context.Context, id string) (string, error)
	getDownloadUrlFn func(ctx context.Context, id string, filename string) (string, error)
	deleteObjectFn   func(ctx context.Context, id string) error
	deleteObjectsFn  func(ctx context.Context, id []string) error
}

func (objectStore MockObjectStore) NewUploadUrl(ctx context.Context, id string) (string, error) {
	if objectStore.newUploadUrlFn != nil {
		return objectStore.newUploadUrlFn(ctx, id)
	}

	return "http://example.org", nil
}

func (objectStore MockObjectStore) GetDownloadUrl(ctx context.Context, id string, filename string) (string, error) {
	if objectStore.getDownloadUrlFn != nil {
		return objectStore.getDownloadUrlFn(ctx, id, filename)
	}

	return "http://example.org", nil
}

func (objectStore MockObjectStore) DeleteObject(ctx context.Context, id string) error {
	if objectStore.deleteObjectFn != nil {
		return objectStore.deleteObjectFn(ctx, id)
	}

	return nil
}

func (objectStore MockObjectStore) DeleteObjects(ctx context.Context, ids []string) error {
	if objectStore.deleteObjectsFn != nil {
		return objectStore.deleteObjectsFn(ctx, ids)
	}

	return nil
}

func NewMockObjectStore() MockObjectStore {
	return MockObjectStore{}
}

func (objectstore MockObjectStore) WithNewUploadUrl(fn func(ctx context.Context, id string) (string, error)) MockObjectStore {
	objectstore.newUploadUrlFn = fn
	return objectstore
}

func (objectstore MockObjectStore) WithGetDownloadUrl(fn func(ctx context.Context, id string, filename string) (string, error)) MockObjectStore {
	objectstore.getDownloadUrlFn = fn
	return objectstore
}

func (objectstore MockObjectStore) WithDeleteObject(fn func(ctx context.Context, id string) error) MockObjectStore {
	objectstore.deleteObjectFn = fn
	return objectstore
}

func (objectstore MockObjectStore) WithDeleteObjects(fn func(ctx context.Context, ids []string) error) MockObjectStore {
	objectstore.deleteObjectsFn = fn
	return objectstore
}
