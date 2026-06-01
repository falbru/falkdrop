package drop_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/falbru/falkdrop/internal/app/drop"
	objectStoreMock "github.com/falbru/falkdrop/internal/storage/objectstore/mock"
	dropMock "github.com/falbru/falkdrop/internal/storage/repository/mock"
	"github.com/google/uuid"
)

func getService(repository dropMock.MockDropRepository, objectStore objectStoreMock.MockObjectStore) *drop.DropService {
	return drop.NewDropService(
		repository,
		objectStore,
	)
}

func TestCreateDrop(t *testing.T) {
	t.Run("empty resource list", func(t *testing.T) {
		service := getService(dropMock.NewMockDropRepository(), objectStoreMock.NewMockObjectStore())

		_, err := service.CreateDrop([]drop.ResourceId{})

		if err == nil {
			t.Errorf("Expected service to throw error with empty resource list")
		} else if err.Error() != "resourceIds can't be empty" {
			t.Errorf("Expected service to throw empty resourceid error, but got: %v", err.Error())
		}
	})

	t.Run("resource in resourceIds doesn't exist", func(t *testing.T) {
		repo := dropMock.NewMockDropRepository().WithGetResourcesByIds(
			func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error) {
				return []drop.Resource{}, nil
			},
		)

		service := getService(repo, objectStoreMock.NewMockObjectStore())
		resourceId := drop.ResourceId(uuid.New())

		_, err := service.CreateDrop([]drop.ResourceId{resourceId})

		var expectedError drop.ErrResourcesNotFound
		if err == nil {
			t.Errorf("Expected service to throw error")
		} else if !errors.As(err, &expectedError) {
			t.Errorf("Expected service to throw ErrResourcesNotFound, but got error: %v", expectedError.Error())
		} else if len(expectedError.ResourceIds) != 1 || expectedError.ResourceIds[0] != resourceId {
			t.Errorf("Expected service to throw ErrResourcesNotFound with resourceId %v, but got error: %v", resourceId, expectedError.Error())
		}
	})

	t.Run("resource in resourceIds is already part of another drop", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		existingDropId := drop.DropId("xxxxx")
		resourceName := "myfile.txt"

		repo := dropMock.NewMockDropRepository().WithGetResourcesByIds(
			func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error) {
				return []drop.Resource{
					{
						Id:     resourceId,
						Type:   drop.FileResource,
						Name:   &resourceName,
						DropId: &existingDropId,
					},
				}, nil
			},
		)

		service := getService(repo, objectStoreMock.NewMockObjectStore())
		_, err := service.CreateDrop([]drop.ResourceId{resourceId})

		var expectedError drop.ErrResourceAlreadyBelongsToDrop
		if err == nil {
			t.Errorf("Expected service to throw error")
		} else if !errors.As(err, &expectedError) {
			t.Errorf("Expected service to throw ErrResourceAlreadyBelongsToDrop, but got error: %v", err.Error())
		} else if expectedError.DropId != existingDropId || expectedError.ResourceId != resourceId {
			t.Errorf("Expected ErrResourceAlreadyBelongsToDrop to have dropId = %v and resourceId = %v, but got: dropId = %v and resourceId = %v", existingDropId, resourceId, expectedError.DropId, expectedError.ResourceId)
		}
	})

	t.Run("success", func(t *testing.T) {
		resourceId := drop.ResourceId(uuid.New())
		resourceType := drop.FileResource
		resourceName := "myfile.txt"
		resourceDownloadUrl := "http://example.org/download"

		objectStore := objectStoreMock.NewMockObjectStore().WithGetDownloadUrl(func(ctx context.Context, id string, filename string) (string, error) {
			return resourceDownloadUrl, nil
		})

		repo := dropMock.NewMockDropRepository().WithGetResourcesByIds(
			func(ctx context.Context, ids []drop.ResourceId) ([]drop.Resource, error) {
				return []drop.Resource{
					{
						Id:     resourceId,
						Type:   resourceType,
						Name:   &resourceName,
						DropId: nil,
					},
				}, nil
			},
		)

		service := getService(repo, objectStore)
		drop, err := service.CreateDrop([]drop.ResourceId{resourceId})

		if err != nil {
			t.Fatalf("Expected no error when calling CreateDrop, but got: %v", err.Error())
		}

		// TODO check expirationDate

		if len(drop.Resources) != 1 {
			t.Errorf("Expected length of drop.Resources to be 1, but got: %v", len(drop.Resources))
		} else {
			resource := drop.Resources[0]

			if resource.Id != resourceId {
				t.Errorf("Expected resource to have id=%v, but got: %v", resourceId, resource.Id)
			}
			if resource.Type != resourceType {
				t.Errorf("Expected resource to have type=%v, but got: %v", resourceType, resource.Type)
			}
			if resource.Name != &resourceName {
				t.Errorf("Expected resource to have name=%v, but got: %v", resourceName, resource.Name)
			}
			if resource.DownloadUrl != resourceDownloadUrl {
				t.Errorf("Expected download URL to be '%v', but got: %v", resourceDownloadUrl, resource.DownloadUrl)
			}
		}
	})
}

func TestGetDropWithResourceDownloadUrls(t *testing.T) {
	t.Run("drop doesn't exist", func(t *testing.T) {
		dropId := drop.DropId("xxxxx")

		repo := dropMock.NewMockDropRepository().WithGetDropById(
			func(ctx context.Context, id drop.DropId) (*drop.Drop, error) {
				return nil, nil
			},
		)

		service := getService(repo, objectStoreMock.NewMockObjectStore())

		_, err := service.GetDropWithResourceDownloadUrls(dropId)

		var expectedError drop.ErrDropNotFound
		if err == nil {
			t.Errorf("Expected service to throw error")
		} else if !errors.As(err, &expectedError) {
			t.Errorf("Expected service to throw ErrDropNotFound, but got error: %v", err.Error())
		} else if expectedError.DropId != dropId {
			t.Errorf("Expected service to throw ErrDropNotFound with dropId %v, but got dropId: %v", dropId, expectedError.DropId)
		}
	})

	t.Run("success", func(t *testing.T) {
		dropId := drop.DropId("xxxxx")
		resourceId := drop.ResourceId(uuid.New())

		objectStore := objectStoreMock.NewMockObjectStore().WithGetDownloadUrl(func(ctx context.Context, id string, filename string) (string, error) {
			return id, nil
		})

		repo := dropMock.NewMockDropRepository().WithGetDropById(
			func(ctx context.Context, id drop.DropId) (*drop.Drop, error) {
				return &drop.Drop{
					Id:             dropId,
					ExpirationDate: time.Now(),
					Resources: []drop.Resource{
						{
							Id:     resourceId,
							Type:   drop.FileResource,
							Name:   nil,
							DropId: &dropId,
						},
					},
				}, nil
			},
		)

		service := getService(repo, objectStore)

		drop, err := service.GetDropWithResourceDownloadUrls(dropId)

		if err != nil {
			t.Fatalf("Expected no error thrown for GetDropWithResourceDownloadUrls")
		}

		if drop.Id != dropId {
			t.Errorf("Expected returned drop to have dropId %v, but got: %v", dropId, drop.Id)
		}

		if len(drop.Resources) != 1 {
			t.Errorf("Expected returned drop to have 1 resource, but resource length is: %v", len(drop.Resources))
		} else {
			resource := drop.Resources[0]

			if resource.Id != resourceId {
				t.Errorf("Expected resource to have resourceId %v, but got: %v", resourceId, resource.Id)
			}

			if resource.DownloadUrl != resource.Id.String() {
				t.Errorf("Expected resource to have downloadURL %v, but got: %v", resource.Id.String(), resource.DownloadUrl)
			}
		}
	})
}

func TestCreateResourceWithUploadUrl(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		resourceType := drop.FileResource
		resourceName := "myfile.txt"
		uploadUrl := "http://myuploadurl.net"

		objectStore := objectStoreMock.NewMockObjectStore().WithNewUploadUrl(func(ctx context.Context, id string) (string, error) {
			return uploadUrl, nil
		})

		repo := dropMock.NewMockDropRepository().WithCreateResource(
			func(ctx context.Context, id drop.ResourceId, resourceType drop.ResourceType, name *string) error {
				return nil
			},
		)

		service := getService(repo, objectStore)

		resource, err := service.CreateResourceWithUploadUrl(resourceType, &resourceName)

		if err != nil {
			t.Fatalf("Expected no error thrown for CreateResourceWithUploadUrl")
		}

		if resource.Type != resourceType {
			t.Errorf("Expected created resource to have type %v, but got: %v", resourceType, resource.Type)
		}

		if resource.Name == nil || *resource.Name != resourceName {
			t.Errorf("Expected created resource to have name %v, but got: %v", resourceName, resource.Name)
		}

		if resource.DropId != nil {
			t.Errorf("Expected created resource to have no dropId, but got: %v", resource.DropId)
		}

		if resource.UploadUrl != uploadUrl {
			t.Errorf("Expected created resource to have uploadUrl %v, but got: %v", uploadUrl, resource.UploadUrl)
		}
	})
}
