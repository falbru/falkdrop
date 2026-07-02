package mock

import "github.com/falbru/falkdrop/internal/app/drop"

type MockDropService struct {
	createResourceWithUploadUrlFn     func(resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error)
	createDropFn                      func(resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error)
	getDropWithResourceDownloadUrlsFn func(dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error)
}

func NewMockDropService() *MockDropService {
	return &MockDropService{}
}

func (service MockDropService) CreateResourceWithUploadUrl(resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
	if service.createResourceWithUploadUrlFn != nil {
		return service.createResourceWithUploadUrlFn(resourceType, name)
	}

	return nil, nil
}

func (service MockDropService) CreateDrop(resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error) {
	if service.createDropFn != nil {
		return service.createDropFn(resourceIds)
	}

	return nil, nil
}

func (service MockDropService) GetDropWithResourceDownloadUrls(dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error) {
	if service.getDropWithResourceDownloadUrlsFn != nil {
		return service.getDropWithResourceDownloadUrlsFn(dropId)
	}

	return nil, nil
}

func (service MockDropService) WithCreateResourceWithUploadUrl(fn func(resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error)) MockDropService {
	service.createResourceWithUploadUrlFn = fn
	return service
}

func (service MockDropService) WithCreateDrop(fn func(resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error)) MockDropService {
	service.createDropFn = fn
	return service
}

func (service MockDropService) WithGetDropWithResourceDownloadUrls(fn func(dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error)) MockDropService {
	service.getDropWithResourceDownloadUrlsFn = fn
	return service
}
