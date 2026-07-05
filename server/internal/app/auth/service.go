package auth

import (
	"context"

	"github.com/falbru/falkdrop/internal/auth"
)

type AuthService struct {
	provider *auth.AuthProvider
}

func NewAuthService(provider *auth.AuthProvider) *AuthService {
	return &AuthService{
		provider,
	}
}

func (service AuthService) Verify(ctx context.Context, token string) error {
	return service.provider.Verify(ctx, token)
}
