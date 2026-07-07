package auth

import (
	"context"

	"github.com/falbru/falkdrop/internal/auth"
)

type AuthService interface {
	Verify(ctx context.Context, token string) error
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
