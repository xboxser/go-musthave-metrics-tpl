package handler

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configSever struct {
	PortSever       string `env:"ADDRESS"`
	IntervalSave    int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"FILE_STORAGE_PATH"`
}

func newConfigServer() *configSever {
	var cfg configSever
	_ = env.Parse(&cfg)

	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	port := agentFlags.String("a", "localhost:8080", "port server")
	intervalSave := agentFlags.Int("i", 300, "time interval save")
	fileStoragePath := agentFlags.String("f", "jsonBD.json", "time interval save")
	restore := agentFlags.Bool("r", true, "time interval save")
	agentFlags.Parse(os.Args[1:])
	if cfg.PortSever == "" {
		cfg.PortSever = *port
	}

	if cfg.IntervalSave == 0 {
		cfg.IntervalSave = *intervalSave
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	if !cfg.Restore {
		cfg.Restore = *restore
	}

	return &cfg
}
