package drop

import (
	"fmt"
	"time"
	"math/rand"
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
func genDropId() string {
	id := make([]rune, 5)

	for i := range id {
		id[i] = dropIdRunes[rand.Intn(len(dropIdRunes))]
	}

	return string(id)
}

func (service DropService) CreateDrop(resourceTypes []ResourceType) error {
	id := genDropId() // TODO check valid id
	expirationDate := time.Now().Add(time.Hour * 24 * 30)

	resources := make([]Resource, len(resourceTypes))
	for i, resourceType := range resourceTypes {
		resources[i] = Resource{
			Link: fmt.Sprintf("TMP-LINK-%v-%v", id, i), // TODO gen link from AWS S3
			Type: resourceType,
		}
	}

	err := service.repository.CreateDrop(id, expirationDate, resources)
	return err
}
