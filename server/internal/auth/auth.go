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

func (auth AuthProvider) HasRole(ctx context.Context, token string, role string) (bool, error) {
	idToken, err := auth.verifier.Verify(ctx, token)
	if err != nil {
		return false, err
	}

	var claims struct {
		ResourceAccess map[string]struct {
			Roles []string `json:"roles"`
		} `json:"resource_access"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return false, err
	}

	if clientRoles, ok := claims.ResourceAccess[auth.clientId]; ok {
		for _, r := range clientRoles.Roles {
			if r == role {
				return true, nil
			}
		}
	}

	return false, nil
}
