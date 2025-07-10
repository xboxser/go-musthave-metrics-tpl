package handler

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configSever struct {
	PortSever string `env:"ADDRESS"`
}

func newConfigServer() *configSever {
	var cfg configSever
	err := env.Parse(&cfg)

	if err != nil || cfg.PortSever == "" {
		agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
		port := agentFlags.String("a", "localhost:8080", "port server")
		agentFlags.Parse(os.Args[1:])
		cfg.PortSever = *port
	}

	return &cfg
}
