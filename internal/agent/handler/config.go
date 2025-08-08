package agent

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configAgent struct {
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	URL            string `env:"ADDRESS"`
}

func newConfigAgent() *configAgent {
	var cfg configAgent
	_ = env.Parse(&cfg)

	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	url := agentFlags.String("a", "localhost:8080", "port server")
	pollInterval := agentFlags.Int("p", 2, "The interval for building metrics")
	reportInterval := agentFlags.Int("r", 10, "The interval for sending data to the server")
	agentFlags.Parse(os.Args[1:])

	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}
	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}
	if cfg.URL == "" {
		cfg.URL = *url
	}

	return &cfg
}
