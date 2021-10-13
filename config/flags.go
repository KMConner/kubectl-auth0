package config

import (
	"flag"
	"io"
)

func LoadCmdArgs(args []string, console io.Writer) (*CmdLine, error) {
	var conf CmdLine
	flagSet := flag.NewFlagSet("", flag.ContinueOnError)
	flagSet.SetOutput(console)
	flagSet.StringVar(&conf.ContextName, "context", "", "Specify cluster name to sign in")
	flagSet.StringVar(&conf.IdpUrl, "idp-issuer-url", "", "OIDC IDP url")
	flagSet.StringVar(&conf.ClientId, "client-id", "", "Client id")
	err := flagSet.Parse(args)
	return &conf, err
}
