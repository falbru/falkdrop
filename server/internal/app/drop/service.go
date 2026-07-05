package drop

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/falbru/falkdrop/internal/app/auth"
	"github.com/falbru/falkdrop/internal/storage/objectstore"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type DropService interface {
	CreateResourceWithUploadUrl(ctx context.Context, resourceType ResourceType, name *string) (*ResourceWithUploadUrl, error)
	CreateDrop(ctx context.Context, resourceIds []ResourceId) (*DropWithResourceDownloadUrls, error)
	GetDropWithResourceDownloadUrls(ctx context.Context, dropId DropId) (*DropWithResourceDownloadUrls, error)
}

type dropServiceImpl struct {
	repository  DropRepository
	objectStore objectstore.ObjectStore
	authService *auth.AuthService
}

func NewDropService(repository DropRepository, objectStore objectstore.ObjectStore, authService *auth.AuthService) DropService {
	return &dropServiceImpl{
		repository,
		objectStore,
		authService,
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

func (service dropServiceImpl) genUniqueDropId(ctx context.Context) (DropId, error) {
	id := genDropId()

	var err error
	for {
		unique, err := service.repository.IsUniqueDropId(ctx, id)

		if err != nil {
			return "", err
		}

		if !unique {
			id = genDropId()
		} else {
			break
		}
	}

	return id, err
}

func (service dropServiceImpl) CreateResourceWithUploadUrl(ctx context.Context, resourceType ResourceType, name *string) (*ResourceWithUploadUrl, error) {
	token := ctx.Value("token").(string)

	if err := service.authService.Verify(ctx, token); err != nil {
		return nil, err
	}

	id := ResourceId(uuid.New())

	err := service.repository.CreateResource(ctx, id, resourceType, name)
	if err != nil {
		return nil, err
	}

	uploadUrl, err := service.objectStore.NewUploadUrl(ctx, id.String())
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

func (service dropServiceImpl) CreateDrop(ctx context.Context, resourceIds []ResourceId) (*DropWithResourceDownloadUrls, error) {
	token := ctx.Value("token").(string)

	if err := service.authService.Verify(ctx, token); err != nil {
		return nil, err
	}

	if len(resourceIds) == 0 {
		return nil, errors.New("resourceIds can't be empty")
	}

	id, err := service.genUniqueDropId(ctx)
	if err != nil {
		return nil, err
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	// TODO verify if resources are uploaded and not empty?

	resources, err := service.repository.GetResourcesByIds(ctx, resourceIds)
	if err != nil {
		return nil, err
	}

	if len(resources) < len(resourceIds) {
		idFound := make(map[ResourceId]bool, len(resources))
		for _, resource := range resources {
			idFound[resource.Id] = true
		}

		var invalidIds []ResourceId
		for _, resourceId := range resourceIds {
			if !idFound[resourceId] {
				invalidIds = append(invalidIds, resourceId)
			}
		}

		return nil, ErrResourcesNotFound{invalidIds}
	}

	for _, resource := range resources {
		if resource.DropId != nil {
			return nil, ErrResourceAlreadyBelongsToDrop{resource.Id, *resource.DropId}
		}
	}

	err = service.repository.CreateDrop(ctx, id, expirationDate, resourceIds)
	if err != nil {
		return nil, err
	}

	resourcesWithDownloadUrls, err := service.withDownloadUrls(ctx, resources)
	if err != nil {
		return nil, err
	}

	return &DropWithResourceDownloadUrls{
		Id:             id,
		ExpirationDate: expirationDate,
		Resources:      resourcesWithDownloadUrls,
	}, nil
}

func (service dropServiceImpl) GetDropWithResourceDownloadUrls(ctx context.Context, dropId DropId) (*DropWithResourceDownloadUrls, error) {
	drop, err := service.repository.GetDropById(ctx, dropId)
	if err != nil {
		return nil, err
	}

	if drop == nil {
		return nil, ErrDropNotFound{dropId}
	}

	resourcesWithDownloadUrls, err := service.withDownloadUrls(ctx, drop.Resources)
	if err != nil {
		return nil, err
	}

	return &DropWithResourceDownloadUrls{
		Id:             drop.Id,
		ExpirationDate: drop.ExpirationDate,
		Resources:      resourcesWithDownloadUrls,
	}, err
}

func (service dropServiceImpl) withDownloadUrls(ctx context.Context, resources []Resource) ([]ResourceWithDownloadUrl, error) {
	resourcesWithDownloadUrls := make([]ResourceWithDownloadUrl, len(resources))

	var wg errgroup.Group
	for i, resource := range resources {
		wg.Go(func() error {
			resourceName := resource.Id.String()
			if resource.Name != nil {
				resourceName = *resource.Name
			}

			resourceWithDownloadUrl, err := service.objectStore.GetDownloadUrl(ctx, resource.Id.String(), resourceName)

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
