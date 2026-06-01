package drop

import (
	"fmt"
	"strings"
)

type ErrDropNotFound struct {
	DropId DropId
}

func (e ErrDropNotFound) Error() string {
	return fmt.Sprintf("Drop with ID %v not found", e.DropId)
}

type ErrResourcesNotFound struct {
	ResourceIds []ResourceId
}

func (e ErrResourcesNotFound) Error() string {
	resourceIdsStr := make([]string, len(e.ResourceIds))

	for i, resourceId := range e.ResourceIds {
		resourceIdsStr[i] = resourceId.String()
	}

	return fmt.Sprintf("Resources with IDs not found: %s", strings.Join(resourceIdsStr, ", "))
}

type ErrResourceAlreadyBelongsToDrop struct {
	ResourceId ResourceId
	DropId     DropId
}

func (e ErrResourceAlreadyBelongsToDrop) Error() string {
	return fmt.Sprintf("Resource %v already belongs to drop %v", e.ResourceId, e.DropId)
}
