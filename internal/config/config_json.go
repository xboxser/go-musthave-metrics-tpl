package config

import (
	"encoding/json"
	"metrics/internal/storage"
)

type ConfigServerJSON struct {
	Address       string `json:"address"`
	Restore       bool   `json:"restore"`
	StoreInterval int    `json:"store_interval"`
	StoreFile     string `json:"store_file"`
	Database      string `json:"database_dsn"`
	CryptoKey     string `json:"crypto_key"`
}

func NewConfigServerJSON(fileName string) *ConfigServerJSON {
	config := &ConfigServerJSON{}
	storageJSON := storage.NewReadJSON(fileName)
	configBytes, err := storageJSON.GetConfigJSON()
	if err != nil {
		return config
	}
	err = json.Unmarshal(configBytes, config)
	if err != nil {
		return config
	}

	return config
}
