package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewConfigServer - Тестируем конфиг сервера
func TestNewConfigServer(t *testing.T) {

	oldArgs := os.Args
	oldAddress := os.Getenv("ADDRESS")
	oldIntervalSave := os.Getenv("STORE_INTERVAL")
	oldFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	oldDateBaseDSN := os.Getenv("DATABASE_DSN")
	oldKEY := os.Getenv("KEY")
	oldRestore := os.Getenv("RESTORE")
	oldAuditFile := os.Getenv("AUDIT_FILE")
	oldAuditURL := os.Getenv("AUDIT_URL")
	oldCryptoKeyPath := os.Getenv("CRYPTO_KEY")

	// сохраняем старые значения, чтобы потом восстановить после прохождения тестов
	defer func() {
		os.Args = oldArgs
		os.Setenv("ADDRESS", oldAddress)
		os.Setenv("STORE_INTERVAL", oldIntervalSave)
		os.Setenv("FILE_STORAGE_PATH", oldFileStoragePath)
		os.Setenv("DATABASE_DSN", oldDateBaseDSN)
		os.Setenv("KEY", oldKEY)
		os.Setenv("RESTORE", oldRestore)
		os.Setenv("AUDIT_FILE", oldAuditFile)
		os.Setenv("AUDIT_URL", oldAuditURL)
		os.Setenv("CRYPTO_KEY", oldCryptoKeyPath)
	}()

	// Проверяем на получение дефолтных значений
	t.Run("checking the receipt of default parameters", func(t *testing.T) {

		os.Unsetenv("ADDRESS")
		os.Unsetenv("STORE_INTERVAL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("RESTORE")
		os.Unsetenv("CRYPTO_KEY")

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigServer()
		require.NotEmpty(t, cfg.Address)
		require.NotEmpty(t, cfg.IntervalSave)
		require.NotEmpty(t, cfg.FileStoragePath)
		require.NotEmpty(t, cfg.Restore)
		require.NotEmpty(t, cfg.CryptoKeyPrivatePath)
		require.Empty(t, cfg.KEY)
		require.Empty(t, cfg.DateBaseDSN)
		require.Empty(t, cfg.AuditFile)
		require.Empty(t, cfg.AuditURL)

	})

	// Тестируем, что при наличии переменных окружения, значения они корректно попадают в конфиг
	t.Run("checking set env parameters", func(t *testing.T) {

		URL := "localhost:8080"
		interval := 10
		KEY := "1234567890"
		cryptoKeyPath := "/private/"

		os.Setenv("ADDRESS", URL)
		os.Setenv("STORE_INTERVAL", strconv.Itoa(interval))
		os.Setenv("FILE_STORAGE_PATH", "FILE_STORAGE_PATH")
		os.Setenv("DATABASE_DSN", "DATABASE_DSN")
		os.Setenv("KEY", KEY)
		os.Setenv("CRYPTO_KEY", cryptoKeyPath)
		os.Setenv("RESTORE", "true")
		os.Setenv("AUDIT_FILE", "AUDIT_FILE")
		os.Setenv("AUDIT_URL", "AUDIT_URL")

		// Устанавливаем значения, чтобы не было ошибок у парсера флагов
		os.Args = []string{"program"}

		cfg := NewConfigServer()
		require.Equal(t, cfg.Address, URL)
		require.Equal(t, cfg.IntervalSave, interval)
		require.Equal(t, cfg.FileStoragePath, "FILE_STORAGE_PATH")
		require.Equal(t, cfg.Restore, true)
		require.Equal(t, cfg.KEY, KEY)
		require.Equal(t, cfg.CryptoKeyPrivatePath, cryptoKeyPath)
		require.Equal(t, cfg.DateBaseDSN, "DATABASE_DSN")
		require.Equal(t, cfg.AuditFile, "AUDIT_FILE")
		require.Equal(t, cfg.AuditURL, "AUDIT_URL")
	})
}
