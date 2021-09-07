package config

import "errors"

type Oidc struct {
	ClientId     string
	Token        string
	IdpUrl       string
	RefreshToken string
}

func LoadOidcConfig(conf *Oidc, configs map[string]string) error {
	clientId, ok := configs["client-id"]
	if !ok {
		return errors.New("key open-id not found")
	}
	conf.ClientId = clientId

	token, ok := configs["id-token"]
	if ok {
		conf.Token = token
	}

	idpUrl, ok := configs["idp-issuer-url"]
	if !ok {
		return errors.New("key idp-issuer-url not found")
	}
	conf.IdpUrl = idpUrl

	refreshToken, ok := configs["refresh-token"]
	if ok {
		conf.RefreshToken = refreshToken
	}

	return nil
}
