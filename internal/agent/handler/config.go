package agent

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configAgent struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	Url            string `env:"ADDRESS"`
}

func newCongigAgent() *configAgent {
	var cfg configAgent
	_ = env.Parse(&cfg)
	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	if cfg.PollInterval == 0 {
		cfg.PollInterval = *agentFlags.Int("p", 2, "The interval for building metrics")
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *agentFlags.Int("r", 10, "The interval for sending data to the server")
	}
	if cfg.Url == "" {
		cfg.Url = *agentFlags.String("a", "localhost:8080", "port server")
	}

	agentFlags.Parse(os.Args[1:])

	return &cfg
}
