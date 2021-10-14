package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/KMConner/kubectl-auth0/config"
	"github.com/KMConner/kubectl-auth0/web"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"

	"github.com/coreos/go-oidc"
)

func ProcessSignIn(oidcReq *config.OidcRequest) (*config.Oidc, error) {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, oidcReq.IdpUrl)
	if err != nil {
		return nil, err
	}

	cnf := &oauth2.Config{
		ClientID:    oidcReq.ClientId,
		Endpoint:    provider.Endpoint(),
		RedirectURL: "http://localhost:8088/callback",
		Scopes:      []string{oidc.ScopeOpenID, oidc.ScopeOfflineAccess, "profile", "email"},
	}

	b := make([]byte, 32)
	_, err = rand.Read(b)
	if err != nil {
		return nil, err
	}
	state := base64.StdEncoding.EncodeToString(b)

	authUrl := cnf.AuthCodeURL(state)
	fmt.Printf("Opening URL %s\n", authUrl)
	err = browser.OpenURL(authUrl)
	if err != nil {
		fmt.Printf("failed to open browser: %+v\nPlease open manually\n", err)
	}

	result, err := web.WaitCallback(func(m map[string][]string) (*web.LoginResult, error) {
		return validateLogin(provider, cnf, state, m)
	})

	if err != nil {
		return nil, err
	}

	return &config.Oidc{
		ClientId:     oidcReq.ClientId,
		Token:        result.Token,
		IdpUrl:       oidcReq.IdpUrl,
		RefreshToken: result.RefreshToken,
	}, nil
}

func validateLogin(provider *oidc.Provider, config *oauth2.Config, state string, query map[string][]string) (*web.LoginResult, error) {
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

	rawIdToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("invalid id_token")
	}

	_, err = provider.Verifier(&oidc.Config{ClientID: config.ClientID}).Verify(context.TODO(), rawIdToken)
	if err != nil {
		return nil, err
	}

	return &web.LoginResult{
		Token:        rawIdToken,
		RefreshToken: token.RefreshToken,
	}, nil
}
