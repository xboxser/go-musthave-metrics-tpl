package storage

import (
	"os"
	"testing"

	models "metrics/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestNewFileJSON(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test_metrics.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Закрываем файл, чтобы NewFileJSON мог его открыть
	tmpFile.Close()

	// Тестируем создание FileJSON
	fileJSON, err := NewFileJSON(tmpFile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, fileJSON)
	assert.NotNil(t, fileJSON.file)
	assert.NotNil(t, fileJSON.encoder)
	assert.NotNil(t, fileJSON.decoder)

	err = fileJSON.Close()
	assert.NoError(t, err)
}

func TestReadEmptyFile(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test_empty_metrics.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Создаем FileJSON
	fileJSON, err := NewFileJSON(tmpFile.Name())
	assert.NoError(t, err)
	defer fileJSON.Close()

	// Читаем из пустого файла
	metrics, err := fileJSON.Read()
	assert.NoError(t, err)
	assert.Empty(t, *metrics)
}

func TestSaveAndRead(t *testing.T) {
	// Создаем временный файл для теста
	tmpFile, err := os.CreateTemp("", "test_metrics.json")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Создаем FileJSON
	fileJSON, err := NewFileJSON(tmpFile.Name())
	assert.NoError(t, err)
	defer fileJSON.Close()

	delta := int64(42)
	value := float64(3.14)
	testMetrics := []models.Metrics{
		{
			ID:    "test_counter",
			MType: "counter",
			Delta: &delta,
		},
		{
			ID:    "test_gauge",
			MType: "gauge",
			Value: &value,
		},
	}

	// Тестируем сохранение
	err = fileJSON.Save(testMetrics)
	assert.NoError(t, err)

	// TODO доделать чтение из файла
}
