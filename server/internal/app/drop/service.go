package drop

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/falbru/falkdrop/internal/storage/objectstore"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type DropService struct {
	repository  DropRepository
	objectStore objectstore.ObjectStore
}

func NewDropService(repository DropRepository, objectStore objectstore.ObjectStore) *DropService {
	return &DropService{
		repository,
		objectStore,
	}
}

var dropIdRunes = []rune("abcdefgijklmnoprstuvwxyz1234567890")

func genDropId() DropId {
	id := make([]rune, 5)

	for i := range id {
		id[i] = dropIdRunes[rand.Intn(len(dropIdRunes))]
	}

	return DropId(id)
}

func (service DropService) genUniqueDropId() (DropId, error) {
	id := genDropId()

	var err error
	for {
		if exists, err := service.repository.IsUniqueDropId(context.Background(), id); exists && err == nil {
			id = genDropId()
		} else {
			break
		}
	}

	return id, err
}

func (service DropService) CreateResourceWithUploadUrl(resourceType ResourceType, name *string) (*ResourceWithUploadUrl, error) {
	id := ResourceId(uuid.New())

	err := service.repository.CreateResource(context.Background(), id, resourceType, name)
	if err != nil {
		return nil, err
	}

	uploadUrl, err := service.objectStore.NewUploadUrl(context.Background(), id.String())
	if err != nil {
		return nil, err
	}

	return &ResourceWithUploadUrl{
		Resource: Resource{
			Id:   id,
			Type: resourceType,
			Name: name,
		},
		UploadUrl: uploadUrl,
	}, nil
}

func (service DropService) CreateDrop(resourceIds []ResourceId) (*DropWithResourceDownloadUrls, error) {
	if len(resourceIds) == 0 {
		return nil, errors.New("resourceIds can't be empty")
	}

	id, err := service.genUniqueDropId()
	if err != nil {
		return nil, err
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	// TODO verify if resources are uploaded and not empty?

	resources, err := service.repository.GetResourcesByIds(context.Background(), resourceIds)
	if err != nil {
		return nil, err
	}

	for _, resource := range resources {
		if resource.DropId != nil {
			return nil, errors.New(fmt.Sprintf("Resource '%v' already part of drop '%v'", resource.Id, resource.DropId))
		}
	}

	err = service.repository.CreateDrop(context.Background(), id, expirationDate, resourceIds)
	if err != nil {
		return nil, err
	}

	resourcesWithDownloadUrls, err := service.withDownloadUrls(resources)
	if err != nil {
		return nil, err
	}

	return &DropWithResourceDownloadUrls{
		Id:             id,
		ExpirationDate: expirationDate,
		Resources:      resourcesWithDownloadUrls,
	}, nil
}

func (service DropService) GetDropWithResourceDownloadUrls(dropId DropId) (*DropWithResourceDownloadUrls, error) {
	drop, err := service.repository.GetDropById(context.Background(), dropId)
	if err != nil {
		return nil, err
	}

	resources, err := service.repository.GetResourcesByDropId(context.Background(), dropId)
	if err != nil {
		return nil, err
	}

	resourcesWithDownloadUrls, err := service.withDownloadUrls(resources)
	if err != nil {
		return nil, err
	}

	return &DropWithResourceDownloadUrls{
		Id:             drop.Id,
		ExpirationDate: drop.ExpirationDate,
		Resources:      resourcesWithDownloadUrls,
	}, err
}

func (service DropService) withDownloadUrls(resources []Resource) ([]ResourceWithDownloadUrl, error) {
	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrl, len(resources))

	var wg errgroup.Group
	for i, resource := range resources {
		wg.Go(func() error {
			resourceName := resource.Id.String()
			if resource.Name != nil {
				resourceName = *resource.Name
			}

			resourceWithDownloadUrl, err := service.objectStore.GetDownloadUrl(context.Background(), resource.Id.String(), resourceName)

			if err != nil {
				return err
			}

			resourcesWithDownloadUrls[i] = ResourceWithDownloadUrl{
				Resource: Resource{
					Id:   resource.Id,
					Type: resource.Type,
					Name: resource.Name,
				},
				DownloadUrl: resourceWithDownloadUrl,
			}

			return nil
		})
	}

	err := wg.Wait()
	if err != nil {
		return []ResourceWithDownloadUrl{}, err
	}

	return resourcesWithDownloadUrls, nil

}
