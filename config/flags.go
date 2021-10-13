package config

import (
	"flag"
	"io"
)

func LoadCmdArgs(args []string, console io.Writer, conf *Config) error {
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.SetOutput(console)
	flagSet.StringVar(&conf.ContextName, "context", "", "Specify cluster name to sign in")
	flagSet.StringVar(&conf.OidcConfig.IdpUrl, "idp-issuer-url", "", "OIDC IDP url")
	flagSet.StringVar(&conf.OidcConfig.ClientId, "client-id", "", "Client id")
	return flagSet.Parse(args)
}
