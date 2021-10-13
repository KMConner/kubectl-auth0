package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/KMConner/kubectl-auth0/config"
	"github.com/KMConner/kubectl-auth0/oauth"
)

func main() {
	var conf config.Config
	err := config.LoadCmdArgs(os.Args[1:], os.Stdout, &conf)
	if err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			log.Fatal(err)
		}
		return
	}

	err = conf.LoadAndValidate()
	if err != nil {
		log.Fatal(err)
	}

	err = oauth.ProcessSignIn(&conf.OidcConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = conf.Save()
	if err != nil {
		log.Fatal(err)
	}
}
