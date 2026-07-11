package mock

import (
	"context"

	"github.com/falbru/falkdrop/internal/app/auth"
)

type MockAuthService struct {
	verifyFn  func(ctx context.Context, token string) error
	hasRoleFn func(ctx context.Context, token string, role auth.AuthRole) (bool, error)
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

func (service MockAuthService) HasRole(ctx context.Context, token string, role auth.AuthRole) (bool, error) {
	if service.hasRoleFn != nil {
		return service.hasRoleFn(ctx, token, role)
	}

	return true, nil
}

func (service MockAuthService) WithVerify(fn func(ctx context.Context, token string) error) *MockAuthService {
	service.verifyFn = fn
	return &service
}

func (service MockAuthService) WithHasRole(fn func(ctx context.Context, token string, role auth.AuthRole) (bool, error)) *MockAuthService {
	service.hasRoleFn = fn
	return &service
}
