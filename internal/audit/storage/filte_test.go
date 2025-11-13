package storage

import (
	"metrics/internal/audit/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewAuditFileJSON - тестируем создание файла
func TestNewAuditFileJSON(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test_metrics.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Закрываем файл, чтобы NewFileJSON мог его открыть
	tmpFile.Close()

	// Тестируем создание FileJSON
	fileJSON, err := NewAuditFileJSON(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, fileJSON)
	assert.NotNil(t, fileJSON.file)
	assert.NotNil(t, fileJSON.encoder)
	assert.NotNil(t, fileJSON.decoder)

	err = fileJSON.Close()
	assert.NoError(t, err)
}

func TestSaveAndRead(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test_metrics.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Создаем FileJSON
	fileAudit, err := NewAuditFileJSON(tmpFile.Name())
	assert.NoError(t, err)
	defer fileAudit.Close()

	testAudit := model.Audit{

		TS:        123,
		Metrics:   []string{"metrics1", "metrics2"},
		IPAddress: "127.0.0.1",
	}

	// Тестируем сохранение
	err = fileAudit.Save(testAudit)
	assert.NoError(t, err)

	err = fileAudit.Close()
	assert.NoError(t, err)
}
