package k8s

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func LoadExisting() (*api.Config, error) {
	pathOpts := clientcmd.NewDefaultPathOptions()
	return pathOpts.GetStartingConfig()
}

func SaveAuth(authInfo *api.AuthInfo, name string) error {
	existing, err := LoadExisting()
	if err != nil {
		return err
	}
	existing.AuthInfos[name] = authInfo
	return clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), *existing, true)
}
