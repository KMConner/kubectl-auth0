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

	oidcReq, username, err := config.LoadOidcConfig(cmdline, k8sConf)
	if err != nil {
		log.Fatal(err)
	}
	if username == "" {
		username = oidcReq.GenerateUsername()
	}

	oidcResult, err := oauth.ProcessSignIn(oidcReq)
	if err != nil {
		log.Fatal(err)
	}
	authInfo := oidcResult.ToAuthInfo()
	err = k8s.SaveAuth(authInfo, username)
	if err != nil {
		log.Fatal(err)
	}
}
