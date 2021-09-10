package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/KMConner/kubectl-auth0/config"
	"github.com/KMConner/kubectl-auth0/web"
	"golang.org/x/oauth2"

	oidc "github.com/coreos/go-oidc"
)

func ProcessSignIn(conf *config.Oidc) error {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, conf.IdpUrl)
	if err != nil {
		return err
	}

	cnf := &oauth2.Config{
		ClientID:    conf.ClientId,
		Endpoint:    provider.Endpoint(),
		RedirectURL: "http://localhost:8088/callback",
		Scopes:      []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"},
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return err
	}
	state := base64.StdEncoding.EncodeToString(b)

	authUrl := cnf.AuthCodeURL(state)
	println(authUrl)

	result, err := web.WaitCallback(func(m map[string][]string) (*web.LoginResult, error) {
		return validateLogin(cnf, state, m)
	})

	if err != nil {
		return err
	}

	conf.Token = result.Token
	conf.RefreshToken = result.RefreshToken
	return nil
}

func validateLogin(config *oauth2.Config, state string, query map[string][]string) (*web.LoginResult, error) {
	codes, ok := query["code"]
	if !ok || len(codes) != 1 {
		return nil, errors.New("invalid code param")
	}
	code := codes[0]

	states, ok := query["state"]
	if !ok || len(states) != 1 {
		return nil, errors.New("invalid state param")
	}
	if state != states[0] {
		return nil, errors.New("state query is wrong")
	}

	token, err := config.Exchange(context.TODO(), code)
	if err != nil {
		return nil, err
	}

	return &web.LoginResult{
		Token:        token.AccessToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
