package config

import (
	"github.com/KMConner/kubectl-auth0/k8s"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Config struct {
	ContextName string
	UserName    string
	OidcConfig  Oidc
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
