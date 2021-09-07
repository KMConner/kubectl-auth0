package main

import (
	"flag"
	"log"

	"github.com/KMConner/kubectl-auth0/config"
)

func main() {
	var conf config.Config
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.StringVar(&conf.ContextName, "context", "", "Specify cluster name to sign in")
	flagSet.StringVar(&conf.OidcConfig.IdpUrl, "idp-issuer-url", "", "OIDC IDP url")
	flagSet.StringVar(&conf.OidcConfig.ClientId, "client-id", "", "Client id")
	err := conf.LoadAndValidate()

	if err != nil {
		log.Fatal(err)
	}

	println("Hello World!")
}
