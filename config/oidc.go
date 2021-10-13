package config

import (
	"crypto/sha256"
	"encoding/base32"
	"errors"

	"k8s.io/client-go/tools/clientcmd/api"
)

type Oidc struct {
	ClientId     string
	Token        string
	IdpUrl       string
	RefreshToken string
}

func LoadOidcConfig(cmdline *CmdLine, k8sConfig *api.Config) (*Oidc, string, error) {
	if cmdline.ClientId != "" && cmdline.IdpUrl != "" {
		return &Oidc{
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

	oidcConf, err := parseOidcConfig(user.AuthProvider.Config)

	if err != nil {
		return nil, "", err
	}

	return oidcConf, ctx.AuthInfo, nil
}

func parseOidcConfig(configs map[string]string) (*Oidc, error) {
	var conf Oidc
	clientId, ok := configs["client-id"]
	if !ok {
		return nil, errors.New("key open-id not found")
	}
	conf.ClientId = clientId

	token, ok := configs["id-token"]
	if ok {
		conf.Token = token
	}

	idpUrl, ok := configs["idp-issuer-url"]
	if !ok {
		return nil, errors.New("key idp-issuer-url not found")
	}
	conf.IdpUrl = idpUrl

	refreshToken, ok := configs["refresh-token"]
	if ok {
		conf.RefreshToken = refreshToken
	}

	return &conf, nil
}

func (o *Oidc) GenerateUsername() string {
	sha := sha256.New()
	hashed := sha.Sum([]byte(o.IdpUrl + o.ClientId))
	return "auth0-" + base32.HexEncoding.EncodeToString(hashed)
}
