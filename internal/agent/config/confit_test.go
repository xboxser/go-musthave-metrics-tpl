package config

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewConfigAgent - Тестируем конфиг сервера
func TestNewConfigAgent(t *testing.T) {

	oldArgs := os.Args
	oldReportInterval := os.Getenv("REPORT_INTERVAL")
	oldPollInterval := os.Getenv("POLL_INTERVAL")
	oldURL := os.Getenv("ADDRESS")
	oldKEY := os.Getenv("KEY")
	oldRateLimit := os.Getenv("RATE_LIMIT")
	oldCryptoKeyPath := os.Getenv("CRYPTO_KEY")

	// сохраняем старые значения, чтобы потом восстановить после прохождения тестов
	defer func() {
		os.Args = oldArgs
		os.Setenv("REPORT_INTERVAL", oldReportInterval)
		os.Setenv("POLL_INTERVAL", oldPollInterval)
		os.Setenv("ADDRESS", oldURL)
		os.Setenv("KEY", oldKEY)
		os.Setenv("RATE_LIMIT", oldRateLimit)
		os.Setenv("CRYPTO_KEY", oldCryptoKeyPath)
	}()

	// Проверяем на получение дефолтных значений
	t.Run("checking the receipt of default parameters", func(t *testing.T) {

		os.Unsetenv("RUN_ADDRESS")
		os.Unsetenv("DATABASE_URI")
		os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
		os.Unsetenv("ACCRUAL_COUNT_CHAN")
		os.Unsetenv("CRYPTO_KEY")

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigAgent()
		require.NotEmpty(t, cfg.ReportInterval)
		require.NotEmpty(t, cfg.PollInterval)
		require.NotEmpty(t, cfg.URL)
		require.NotEmpty(t, cfg.RateLimit)
		require.NotEmpty(t, cfg.CryptoKeyPath)
		require.Empty(t, cfg.KEY)
	})

	// Тестируем, что при наличии переменных окружения, значения они корректно попадают в конфиг
	t.Run("checking set env parameters", func(t *testing.T) {

		URL := "localhost:8080"
		reportInterval := 10
		KEY := "1234567890"
		rateLimit := 2
		pollInterval := 1
		cryptoKeyPath := "/private/"
		os.Setenv("REPORT_INTERVAL", strconv.Itoa(reportInterval))
		os.Setenv("POLL_INTERVAL", strconv.Itoa(pollInterval))
		os.Setenv("ADDRESS", URL)
		os.Setenv("KEY", KEY)
		os.Setenv("CRYPTO_KEY", cryptoKeyPath)
		os.Setenv("RATE_LIMIT", strconv.Itoa(rateLimit))

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigAgent()
		require.Equal(t, cfg.ReportInterval, reportInterval)
		require.Equal(t, cfg.URL, URL)
		require.Equal(t, cfg.KEY, KEY)
		require.Equal(t, cfg.CryptoKeyPath, cryptoKeyPath)
		require.Equal(t, cfg.RateLimit, rateLimit)
		require.Equal(t, cfg.PollInterval, pollInterval)
	})
}

func TestConfigJSON(t *testing.T) {
	t.Run("checking empty path", func(t *testing.T) {
		cfgDef := ConfigAgent{}
		cfg := ConfigAgent{}
		configJSON(&cfg)
		require.Equal(t, cfg, cfgDef)
	})

	t.Run("checking config", func(t *testing.T) {
		// Создаем временный файл
		tmpFile, err := os.CreateTemp("", "config_*.json")
		require.Empty(t, err)
		defer tmpFile.Close()

		// Создаем пример конфигурации
		cfgDef := ConfigAgentJson{
			Address:        "localhost:8080",
			ReportInterval: 10,
			PollInterval:   5,
			CryptoKey:      "/path/to/key.pem",
		}

		configData, err := json.MarshalIndent(cfgDef, "", "  ")
		require.Empty(t, err)

		_, err = tmpFile.Write(configData)
		require.Empty(t, err)
		defer os.Remove(tmpFile.Name())
		cfg := ConfigAgent{ConfigPath: tmpFile.Name()}
		configJSON(&cfg)

		require.Equal(t, cfgDef.Address, cfg.URL)
		require.Equal(t, cfgDef.CryptoKey, cfg.CryptoKeyPath)
		require.Equal(t, cfgDef.PollInterval, cfg.PollInterval)
		require.Equal(t, cfgDef.ReportInterval, cfg.ReportInterval)
	})
}
