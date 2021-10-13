package config

import (
	"errors"

	"k8s.io/client-go/tools/clientcmd/api"
)

type Oidc struct {
	ClientId     string
	Token        string
	IdpUrl       string
	RefreshToken string
}

func LoadOidcConfig(cmdline *CmdLine, k8sConfig *api.Config) (*OidcRequest, string, error) {
	if cmdline.ClientId != "" && cmdline.IdpUrl != "" {
		return &OidcRequest{
			ClientId: cmdline.ClientId,
			IdpUrl:   cmdline.IdpUrl,
		}, "", nil
	}

	contextName := cmdline.ContextName
	if contextName == "" {
		if k8sConfig.CurrentContext == "" {
			return nil, "", errors.New("context is not specified")
		}
		contextName = k8sConfig.CurrentContext
	}

	ctx, ok := k8sConfig.Contexts[contextName]
	if !ok {
		return nil, "", errors.New("context " + contextName + " not found")
	}

	user, ok := k8sConfig.AuthInfos[ctx.AuthInfo]
	if !ok {
		return nil, "", errors.New("user " + ctx.AuthInfo + " not found")
	}
	if user.AuthProvider == nil || user.AuthProvider.Name != "oidc" {
		return nil, "", errors.New("user " + ctx.AuthInfo + " was found but auth provider is not oidc")
	}

	oidcConf, err := newOidcRequest(user.AuthProvider.Config)

	if err != nil {
		return nil, "", err
	}

	return oidcConf, ctx.AuthInfo, nil
}

func (o *Oidc) ToAuthInfo() *api.AuthInfo {
	return &api.AuthInfo{
		AuthProvider: &api.AuthProviderConfig{
			Name: "oidc",
			Config: map[string]string{
				"idp-issuer-url": o.IdpUrl,
				"client-id":      o.ClientId,
				"id-token":       o.Token,
				"refresh-token":  o.RefreshToken,
			},
		},
	}
}
