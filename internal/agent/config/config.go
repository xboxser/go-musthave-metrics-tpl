package config

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
	ConfigPath     string `env:"CONFIG"`     //  Путь до файла с json конфигом
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
	configPath := agentFlags.String("c", "", "path config file")
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

	if cfg.ConfigPath == "" {
		cfg.ConfigPath = *configPath
	}

	configJSON(&cfg)

	return &cfg
}

// configJSON - читаем данные из configJSON
// данные параметры менее приоритетны чем из командной строки или env
func configJSON(c *ConfigAgent) {
	if c.ConfigPath == "" {
		return
	}
	configJSON := NewConfigAgentJson(c.ConfigPath)

	if c.URL == "" {
		c.URL = configJSON.Address
	}

	if c.CryptoKeyPath == "" {
		c.CryptoKeyPath = configJSON.CryptoKey
	}

	if c.PollInterval == 0 {
		c.PollInterval = configJSON.PollInterval
	}
	if c.ReportInterval == 0 {
		c.ReportInterval = configJSON.ReportInterval
	}
}
