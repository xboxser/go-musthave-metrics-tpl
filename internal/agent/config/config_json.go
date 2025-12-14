package config

import (
	"encoding/json"
	"metrics/internal/storage"
)

type ConfigAgentJSON struct {
	Address        string `json:"address"`
	ReportInterval int    `json:"report_interval"`
	PollInterval   int    `json:"poll_interval"`
	CryptoKey      string `json:"crypto_key"`
}

func NewConfigAgentJSON(fileName string) *ConfigAgentJSON {
	config := &ConfigAgentJSON{}
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
