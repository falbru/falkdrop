package auth

import (
	"context"

	"github.com/falbru/falkdrop/internal/auth"
)

type AuthRole string

const (
	CreateDropRole AuthRole = "drops:create"
)

type AuthService interface {
	Verify(ctx context.Context, token string) error
	HasRole(ctx context.Context, token string, role AuthRole) (bool, error)
}

type authServiceImpl struct {
	provider *auth.AuthProvider
}

func NewAuthService(provider *auth.AuthProvider) AuthService {
	return &authServiceImpl{
		provider,
	}
}

func (service authServiceImpl) Verify(ctx context.Context, token string) error {
	return service.provider.Verify(ctx, token)
}

func (service authServiceImpl) HasRole(ctx context.Context, token string, role AuthRole) (bool, error) {
	return service.provider.HasRole(ctx, token, string(role))
}
