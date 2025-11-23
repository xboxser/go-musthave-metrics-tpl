package handler

import (
	"metrics/internal/config"
	models "metrics/internal/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage(t *testing.T) {
	cfg := &config.ConfigServer{}
	cfg.FileStoragePath = "test.json"

	storage := NewStorageManager(cfg)
	defer storage.Close()
	// Проверяем, что структура создана
	assert.NotNil(t, storage)

	// Проверяем, что поля структуры инициализированы
	assert.NotNil(t, storage.config)
	assert.NotNil(t, storage.file)
	assert.Equal(t, storage.config, cfg)

	// Проверяем, что пути к файлу совпадают
	assert.Equal(t, "test.json", storage.config.FileStoragePath)
}

func TestSaveAndRead(t *testing.T) {
	testFileName := "test.json"

	cfg := &config.ConfigServer{}
	cfg.FileStoragePath = testFileName
	storage := NewStorageManager(cfg)
	defer storage.Close()
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

	err := storage.SaveToFile(testMetrics)
	assert.NoError(t, err)

	metrics := storage.ReadFromFile()
	assert.NoError(t, err)
	assert.NotEmpty(t, metrics)

	assert.Equal(t, testMetrics[0].ID, (metrics)[0].ID)
	assert.Equal(t, testMetrics[1].ID, (metrics)[1].ID)
	assert.Equal(t, testMetrics[0].MType, (metrics)[0].MType)
	assert.Equal(t, testMetrics[1].MType, (metrics)[1].MType)
	assert.Equal(t, testMetrics[0].Delta, (metrics)[0].Delta)
	assert.Equal(t, testMetrics[1].Value, (metrics)[1].Value)

	defer func() {
		if _, err := os.Stat(testFileName); err == nil {
			os.Remove(testFileName)
		}
	}()

}
