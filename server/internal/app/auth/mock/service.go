package mock

import (
	"context"
)

type MockAuthService struct {
	verifyFn func(ctx context.Context, token string) error
}

func NewMockAuthService() *MockAuthService {
	return &MockAuthService{}
}

func (service MockAuthService) Verify(ctx context.Context, token string) error {
	if service.verifyFn != nil {
		return service.verifyFn(ctx, token)
	}

	return nil
}

func (service MockAuthService) WithVerify(fn func(ctx context.Context, token string) error) *MockAuthService {
	service.verifyFn = fn
	return &service
}
