package config

import (
	"crypto/sha256"
	"encoding/base32"
	"errors"
)

type OidcRequest struct {
	ClientId string
	IdpUrl   string
}

func newOidcRequest(params map[string]string) (*OidcRequest, error) {
	var conf OidcRequest
	clientId, ok := params["client-id"]
	if !ok {
		return nil, errors.New("key open-id not found")
	}
	conf.ClientId = clientId

	idpUrl, ok := params["idp-issuer-url"]
	if !ok {
		return nil, errors.New("key idp-issuer-url not found")
	}
	conf.IdpUrl = idpUrl
	return &conf, nil
}

func (o *OidcRequest) GenerateUsername() string {
	sha := sha256.New()
	hashed := sha.Sum([]byte(o.IdpUrl + o.ClientId))
	return "auth0-" + base32.HexEncoding.EncodeToString(hashed)[:5]
}
