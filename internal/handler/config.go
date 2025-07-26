package handler

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configServer struct {
	PortSever       string `env:"ADDRESS"`
	IntervalSave    int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func newConfigServer() *configServer {
	var cfg configServer
	_ = env.Parse(&cfg)

	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	port := serverFlags.String("a", "localhost:8080", "port server")
	intervalSave := serverFlags.Int("i", 300, "time interval save")
	fileStoragePath := serverFlags.String("f", "jsonBD.json", "time interval save")
	restore := serverFlags.Bool("r", true, "read file to start server")
	serverFlags.Parse(os.Args[1:])
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
