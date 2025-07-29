package handler

import (
	"flag"
	"os"

	"github.com/caarlos0/env"
)

type configServer struct {
	Address       string `env:"ADDRESS"`
	IntervalSave    int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DateBaseDSN string `env:"DATABASE_DSN"`
	Restore         bool   `env:"RESTORE"`
}

func newConfigServer() *configServer {
	var cfg configServer
	_ = env.Parse(&cfg)

	serverFlags := flag.NewFlagSet("server", flag.ExitOnError)
	address := serverFlags.String("a", "localhost:8080", "port server")
	intervalSave := serverFlags.Int("i", 300, "time interval save")
	fileStoragePath := serverFlags.String("f", "jsonBD.json", "the path to the file to save the data")
	dateBaseDSN := serverFlags.String("d", "localhost:5432", "host db PostgreSQL")
	restore := serverFlags.Bool("r", true, "read file to start server")
	serverFlags.Parse(os.Args[1:])
	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.IntervalSave == 0 {
		cfg.IntervalSave = *intervalSave
	}

	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}

	if cfg.DateBaseDSN == "" {
		cfg.DateBaseDSN = *dateBaseDSN
	}

	if !cfg.Restore {
		cfg.Restore = *restore
	}

	return &cfg
}
