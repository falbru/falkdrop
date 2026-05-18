package objectstore

type ObjectStore interface {
	NewUploadUrl(id string) (string, error)
}
