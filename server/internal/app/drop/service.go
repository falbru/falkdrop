package drop

import (
	"context"
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

func (service DropService) CreateResourceWithUploadLink(resourceType ResourceType) (ResourceId, string, error) {
	id := ResourceId(uuid.New())

	err := service.repository.CreateResource(context.Background(), id, resourceType)
	if err != nil {
		return ResourceId(uuid.Nil), "", err
	}

	uploadUrl, err := service.objectStore.NewUploadUrl(context.Background(), id.String())
	return id, uploadUrl, err
}

func (service DropService) CreateDrop(resourceIds []ResourceId) (DropId, error) {
	id, err := service.genUniqueDropId()
	if err != nil {
		return "", err
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	// TODO verify if resources are uploaded and not empty?

	return id, service.repository.CreateDrop(context.Background(), id, expirationDate, resourceIds)
}

func (service DropService) GetDropWithResourceLinks(dropId DropId) (*DropWithResourceLinks, error) {
	drop, err := service.repository.GetDropById(context.Background(), dropId)
	if err != nil {
		return nil, err
	}

	resources, err := service.repository.GetResourcesByDropId(context.Background(), dropId)

	resourceLinks := make([]ResourceLink, len(resources))

	var wg errgroup.Group
	for i, resource := range resources {
		wg.Go(func() error {
			resourceLink, err := service.objectStore.GetDownloadUrl(context.Background(), resource.Id.String())

			if err != nil {
				return err
			}

			resourceLinks[i] = ResourceLink{
				Type: resource.Type,
				Link: resourceLink,
			}

			return nil
		})
	}

	err = wg.Wait()

	if err != nil {
		return nil, err
	}

	return &DropWithResourceLinks{
		Id:             drop.Id,
		ExpirationDate: drop.ExpirationDate,
		ResourceLinks:  resourceLinks,
	}, err
}
