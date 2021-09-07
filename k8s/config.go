package k8s

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func LoadExisting() (*api.Config, error) {
	pathOpts := clientcmd.NewDefaultPathOptions()
	return pathOpts.GetStartingConfig()
}
