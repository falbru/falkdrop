package mock

import (
	"context"

	"github.com/falbru/falkdrop/internal/app/auth"
	authMock "github.com/falbru/falkdrop/internal/app/auth/mock"
	"github.com/falbru/falkdrop/internal/app/drop"
)

type MockDropService struct {
	authService                       auth.AuthService
	createResourceWithUploadUrlFn     func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error)
	createDropFn                      func(ctx context.Context, resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error)
	getDropWithResourceDownloadUrlsFn func(ctx context.Context, dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error)
}

func NewMockDropService() *MockDropService {
	return &MockDropService{
		authService: authMock.NewMockAuthService(),
	}
}

func (service MockDropService) CreateResourceWithUploadUrl(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error) {
	if service.createResourceWithUploadUrlFn != nil {
		return service.createResourceWithUploadUrlFn(ctx, resourceType, name)
	}

	return nil, nil
}

func (service MockDropService) CreateDrop(ctx context.Context, resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error) {
	if service.createDropFn != nil {
		return service.createDropFn(ctx, resourceIds)
	}

	return nil, nil
}

func (service MockDropService) GetDropWithResourceDownloadUrls(ctx context.Context, dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error) {
	if service.getDropWithResourceDownloadUrlsFn != nil {
		return service.getDropWithResourceDownloadUrlsFn(ctx, dropId)
	}

	return nil, nil
}

func (service MockDropService) WithAuthService(authService auth.AuthService) *MockDropService {
	service.authService = authService
	return &service
}

func (service MockDropService) WithCreateResourceWithUploadUrl(fn func(ctx context.Context, resourceType drop.ResourceType, name *string) (*drop.ResourceWithUploadUrl, error)) *MockDropService {
	service.createResourceWithUploadUrlFn = fn
	return &service
}

func (service MockDropService) WithCreateDrop(fn func(ctx context.Context, resourceIds []drop.ResourceId) (*drop.DropWithResourceDownloadUrls, error)) *MockDropService {
	service.createDropFn = fn
	return &service
}

func (service MockDropService) WithGetDropWithResourceDownloadUrls(fn func(ctx context.Context, dropId drop.DropId) (*drop.DropWithResourceDownloadUrls, error)) *MockDropService {
	service.getDropWithResourceDownloadUrlsFn = fn
	return &service
}
