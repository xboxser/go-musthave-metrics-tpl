package agent

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/caarlos0/env"
)

// ConfigAgent - конфигурация агента
type ConfigAgent struct {
	ReportInterval int    `env:"REPORT_INTERVAL"` // Интервал отправки данных на сервер, в секундах
	PollInterval   int    `env:"POLL_INTERVAL"`   // Интервал сбора метрик, в секундах
	URL            string `env:"ADDRESS"`         // Адрес получателя метрик
	KEY            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
	CryptoKeyPath  string `env:"CRYPTO_KEY"` //  Путь до файла с публичным ключом
}

func NewConfigAgent() *ConfigAgent {
	var cfg ConfigAgent
	_ = env.Parse(&cfg)

	homePath, _ := os.UserHomeDir()
	homePath = filepath.Join(homePath, "cert.pem")

	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)
	url := agentFlags.String("a", "localhost:8080", "port server")
	pollInterval := agentFlags.Int("p", 2, "The interval for building metrics")
	reportInterval := agentFlags.Int("r", 10, "The interval for sending data to the server")
	key := agentFlags.String("k", "", "specify the encryption key")
	rateLimit := agentFlags.Int("l", 2, "rate limit")
	cryptoKeyPath := agentFlags.String("crypto-key", homePath, "path crypto key")
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
	if cfg.KEY == "" {
		cfg.KEY = *key
	}
	if cfg.RateLimit == 0 {
		cfg.RateLimit = *rateLimit
	}
	if cfg.CryptoKeyPath == "" {
		cfg.CryptoKeyPath = *cryptoKeyPath
	}
	return &cfg
}
