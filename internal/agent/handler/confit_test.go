package agent

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewConfigServer - Тестируем конфиг сервера
func TestNewConfigServer(t *testing.T) {

	oldArgs := os.Args
	oldReportInterval := os.Getenv("REPORT_INTERVAL")
	oldPollInterval := os.Getenv("POLL_INTERVAL")
	oldURL := os.Getenv("ADDRESS")
	oldKEY := os.Getenv("KEY")
	oldRateLimit := os.Getenv("RATE_LIMIT")

	// сохраняем старые значения, чтобы потом восстановить после прохождения тестов
	defer func() {
		os.Args = oldArgs
		os.Setenv("REPORT_INTERVAL", oldReportInterval)
		os.Setenv("POLL_INTERVAL", oldPollInterval)
		os.Setenv("ADDRESS", oldURL)
		os.Setenv("KEY", oldKEY)
		os.Setenv("RATE_LIMIT", oldRateLimit)
	}()

	// Проверяем на получение дефолтных значений
	t.Run("checking the receipt of default parameters", func(t *testing.T) {

		os.Unsetenv("RUN_ADDRESS")
		os.Unsetenv("DATABASE_URI")
		os.Unsetenv("ACCRUAL_SYSTEM_ADDRESS")
		os.Unsetenv("ACCRUAL_COUNT_CHAN")

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigAgent()
		require.NotEmpty(t, cfg.ReportInterval)
		require.NotEmpty(t, cfg.PollInterval)
		require.NotEmpty(t, cfg.URL)
		require.NotEmpty(t, cfg.RateLimit)
		require.Empty(t, cfg.KEY)
	})

	// Тестируем, что при наличии переменных окружения, значения они корректно попадают в конфиг
	t.Run("checking set env parameters", func(t *testing.T) {

		URL := "localhost:8080"
		reportInterval := 10
		KEY := "1234567890"
		rateLimit := 2
		pollInterval := 1

		os.Setenv("REPORT_INTERVAL", strconv.Itoa(reportInterval))
		os.Setenv("POLL_INTERVAL", strconv.Itoa(pollInterval))
		os.Setenv("ADDRESS", URL)
		os.Setenv("KEY", KEY)
		os.Setenv("RATE_LIMIT", strconv.Itoa(rateLimit))

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigAgent()
		require.Equal(t, cfg.ReportInterval, reportInterval)
		require.Equal(t, cfg.URL, URL)
		require.Equal(t, cfg.KEY, KEY)
		require.Equal(t, cfg.RateLimit, rateLimit)
		require.Equal(t, cfg.PollInterval, pollInterval)
	})
}
