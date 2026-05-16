package drop

import (
	"fmt"
	"math/rand"
	"time"
)

type DropService struct {
	repository DropRepository
}

func NewDropService(repository DropRepository) *DropService {
	return &DropService{
		repository,
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

func (service DropService) CreateDrop(resourceTypes []ResourceType) error {
	id, err := service.genUniqueDropId()
	if err != nil {
		return err
	}

	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	resources := make([]Resource, len(resourceTypes))
	for i, resourceType := range resourceTypes {
		resources[i] = Resource{
			Link: fmt.Sprintf("TMP-LINK-%v-%v", id, i), // TODO gen link from AWS S3
			Type: resourceType,
		}
	}

	err = service.repository.CreateDrop(id, expirationDate, resources)
	return err
}
