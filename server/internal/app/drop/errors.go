package drop

import "fmt"

type DropNotFound struct {
	DropId DropId
}

func (e DropNotFound) Error() string {
	return fmt.Sprintf("Drop with ID %v not found", e.DropId)
}
