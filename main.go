package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/KMConner/kubectl-auth0/config"
	"github.com/KMConner/kubectl-auth0/k8s"
	"github.com/KMConner/kubectl-auth0/oauth"
)

func main() {
	cmdline, err := config.LoadCmdArgs(os.Args[1:], os.Stdout)
	if err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			log.Fatal(err)
		}
		return
	}

	k8sConf, err := k8s.LoadExisting()
	if err != nil {
		log.Fatal(err)
	}

	oidcConfig, username, err := config.LoadOidcConfig(cmdline, k8sConf)
	if err != nil {
		log.Fatal(err)
	}
	if username == "" {
		username = oidcConfig.GenerateUsername()
	}

	err = oauth.ProcessSignIn(oidcConfig)
	if err != nil {
		log.Fatal(err)
	}

	conf := &config.Config{
		UserName:   username,
		OidcConfig: *oidcConfig,
	}
	err = conf.Save()
	if err != nil {
		log.Fatal(err)
	}
}
