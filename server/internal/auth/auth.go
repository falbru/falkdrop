package auth

import (
	"context"
	"github.com/coreos/go-oidc/v3/oidc"
)

type AuthProvider struct {
	realmConfigURL string
	clientId       string

	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
}

func NewAuthProvider(RealmConfigURL string, ClientId string) *AuthProvider {
	return &AuthProvider{
		realmConfigURL: RealmConfigURL,
		clientId:       ClientId,
	}
}

func (auth *AuthProvider) Init(ctx context.Context) error {
	provider, err := oidc.NewProvider(ctx, auth.realmConfigURL)
	if err != nil {
		return err
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: auth.clientId,
	})

	auth.provider = provider
	auth.verifier = verifier

	return nil
}

func (auth AuthProvider) Verify(ctx context.Context, token string) error {
	_, err := auth.verifier.Verify(context.Background(), token)

	return err
}
