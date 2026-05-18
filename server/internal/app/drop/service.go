package drop

import (
	"math/rand"
	"time"

	"github.com/falbru/falkdrop/internal/storage/objectstore"
	"github.com/google/uuid"
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
		if exists, err := service.repository.IsUniqueDropId(id); exists && err == nil {
			id = genDropId()
		} else {
			break
		}
	}

	return id, err
}

func (service DropService) CreateResourceWithUploadLink(resourceType ResourceType) (ResourceId, string, error) {
	id := ResourceId(uuid.New())

	err := service.repository.CreateResource(id, resourceType)
	if err != nil {
		return ResourceId(uuid.Nil), "", err
	}

	uploadUrl, err := service.objectStore.NewUploadUrl(id.String())
	return id, uploadUrl, err
}

func (service DropService) CreateDrop(resourceIds []ResourceId) error {
	id, err := service.genUniqueDropId()
	if err != nil {
		return err
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	return service.repository.CreateDrop(id, expirationDate, resourceIds)
}
