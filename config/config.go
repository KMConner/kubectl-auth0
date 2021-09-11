package config

import (
	"errors"
	"fmt"

	"github.com/KMConner/kubectl-auth0/k8s"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Config struct {
	ContextName string
	UserName    string
	OidcConfig  Oidc
}

func (c *Config) LoadAndValidate() error {
	if c.OidcConfig.ClientId != "" && c.OidcConfig.IdpUrl != "" {
		return nil
	}

	existingConfig, err := k8s.LoadExisting()
	if err != nil {
		return err
	}

	if c.ContextName == "" {
		fmt.Printf("default context %s is used\n", existingConfig.CurrentContext)
		c.ContextName = existingConfig.CurrentContext
	}

	ctx, ok := existingConfig.Contexts[c.ContextName]
	if !ok {
		return errors.New("context " + c.ContextName + " not found")
	}
	c.UserName = ctx.AuthInfo
	user, ok := existingConfig.AuthInfos[ctx.AuthInfo]
	if !ok {
		return errors.New("user " + ctx.AuthInfo + " not found")
	}
	if user.AuthProvider == nil || user.AuthProvider.Name != "oidc" {
		return errors.New("user " + ctx.AuthInfo + " was found but auth provider is not oidc")
	}

	err = LoadOidcConfig(&c.OidcConfig, user.AuthProvider.Config)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) Save() error {
	config := &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"idp-issuer-url": c.OidcConfig.IdpUrl,
				"client-id":      c.OidcConfig.ClientId,
				"id-token":       c.OidcConfig.Token,
				"refresh-token":  c.OidcConfig.RefreshToken,
			},
		},
	}
	return k8s.SaveAuth(config, c.UserName)
}
