package service

import (
	models "metrics/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

// _ "net/http/pprof" // подключаем пакет pprof

func TestNewServeService(t *testing.T) {
	model := models.NewMemStorage()
	service := NewServeService(model)

	if service == nil {
		t.Fatal("failed to create service")
	}
	if service.model != model {
		t.Fatal("NewSender returned wrong model")
	}
}

// TestUpdateAndGetCounter - проверка методов Update и GetValue для типа counter
func TestUpdateAndGetCounter(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Тестируем обновление счетчика
	err := service.Update("counter", "test_counter", "42")
	assert.NoError(t, err)

	// Тестируем получение значения счетчика
	value, err := service.GetValue("counter", "test_counter")
	assert.NoError(t, err)
	assert.Equal(t, "42", value)
}

// TestUpdateAndGetGauge - проверка методов Update и GetValue для типа gauge
func TestUpdateAndGetGauge(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Тестируем обновление gauge
	err := service.Update("gauge", "test_gauge", "3.14")
	assert.NoError(t, err)

	// Тестируем получение значения gauge
	value, err := service.GetValue("gauge", "test_gauge")
	assert.NoError(t, err)
	assert.Equal(t, "3.14", value)
}

// TestUpdateError - проверяем получение ошибки при обновлении и получение метрик
func TestUpdateError(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Тестируем обновление не корректного типа
	err := service.Update("randomName", "test_gauge", "3.14")
	assert.Error(t, err)

	// Тестируем получение не корректного типа
	value, err := service.GetValue("randomName", "test_gauge")
	assert.Error(t, err)
	assert.Equal(t, "", value)

	// Тестируем обновление gauge/ Тестируем обновление не корректного значения
	err = service.Update("counter", "test_counter", "3.14")
	assert.Error(t, err)
	err = service.Update("gauge", "test_gauge2", "text")
	assert.Error(t, err)

	value, err = service.GetValue("counter", "test_counter_123")
	assert.Error(t, err)
	assert.Equal(t, "", value)

	value, err = service.GetValue("gauge", "test_gauge_123")
	assert.Error(t, err)
	assert.Equal(t, "", value)
}

// TestUpdateJSONCounter - проверяем UpdateJSON для типа counter
func TestUpdateJSONCounter(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	delta := int64(10)
	metric := &models.Metrics{
		ID:    "test_counter",
		MType: "counter",
		Delta: &delta,
	}

	err := service.UpdateJSON(metric)
	assert.NoError(t, err)

	// Проверяем, что Delta был обновлен после операции
	assert.Equal(t, int64(10), *metric.Delta)
}

// TestUpdateJSONGauge - проверяем UpdateJSON для типа gauge
func TestUpdateJSONGauge(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	value := float64(100)
	metric := &models.Metrics{
		ID:    "test_counter",
		MType: "gauge",
		Value: &value,
	}

	err := service.UpdateJSON(metric)
	assert.NoError(t, err)

	// Проверяем, что Value был обновлен после операции
	assert.Equal(t, value, *metric.Value)
}

// TestUpdateJSONError - проверяем не корректное получение UpdateJSON
func TestUpdateJSONError(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	value := float64(100)
	metric := &models.Metrics{
		ID:    "test_counter",
		MType: "randomName",
		Value: &value,
	}
	// Проверяем не корректный тип MType
	err := service.UpdateJSON(metric)
	assert.Error(t, err)
}

// TestGetValueJSON - проверяет работу метода GetValueJSON
func TestGetValueJSON(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	valueFloat := float64(100)
	valueInt := int64(100)
	metricGauge := &models.Metrics{
		ID:    "test_counter",
		MType: "gauge",
		Value: &valueFloat,
	}

	metricCounter := &models.Metrics{
		ID:    "test_counter",
		MType: "counter",
		Delta: &valueInt,
	}

	err := service.UpdateJSON(metricGauge)
	assert.NoError(t, err)

	err = service.UpdateJSON(metricCounter)
	assert.NoError(t, err)

	// Проверяем получение ошибок
	getMetric := &models.Metrics{
		ID:    "test_counter_123",
		MType: "gauge",
	}

	err = service.GetValueJSON(getMetric)
	assert.Error(t, err)

	getMetric.MType = "counter"
	err = service.GetValueJSON(getMetric)
	assert.Error(t, err)

	getMetric.MType = "gauge"
	err = service.GetValueJSON(getMetric)
	assert.Error(t, err)

	// Проверяем корректные значения
	getMetric.ID = "test_counter"
	err = service.GetValueJSON(getMetric)
	assert.NoError(t, err)
	assert.Equal(t, valueFloat, *getMetric.Value)

	getMetric.MType = "counter"
	err = service.GetValueJSON(getMetric)
	assert.NoError(t, err)
	assert.Equal(t, valueInt, *getMetric.Delta)
}

// TestGetAll - тестируем метод GetAll
func TestGetAll(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Добавляем тестовые данные
	service.Update("counter", "test_counter", "10")
	service.Update("gauge", "test_gauge", "3.14")

	// Получаем все метрики
	allMetrics := service.GetAll()

	// Проверяем, что метрики присутствуют
	assert.Contains(t, allMetrics, "counter test_counter")
	assert.Contains(t, allMetrics, "gauge test_gauge")
	assert.Equal(t, "10", allMetrics["counter test_counter"])
	assert.Equal(t, "3.14", allMetrics["gauge test_gauge"])
}

// TestGetModels - тестируем метод GetModels
func TestGetModels(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Подготовка тестовых данных
	// Добавляем counter
	service.Update("counter", "test_counter", "42")

	// Добавляем gauge
	service.Update("gauge", "test_gauge", "3.14159")

	// Вызываем тестируемую функцию
	models := service.GetModels()

	// Проверяем, что функция вернула непустой слайс
	assert.NotEmpty(t, models)

	// Проверяем, что в результате есть обе метрики
	assert.Len(t, models, 2)

	// Проверяем содержимое метрик
	foundCounter := false
	foundGauge := false

	for _, m := range models {
		if m.ID == "test_counter" && m.MType == "counter" {
			assert.Equal(t, int64(42), *m.Delta)
			foundCounter = true
		}
		if m.ID == "test_gauge" && m.MType == "gauge" {
			assert.Equal(t, 3.14159, *m.Value)
			foundGauge = true
		}
	}

	// Убеждаемся, что обе метрики были найдены
	assert.True(t, foundCounter, "Counter metric not found in GetModels result")
	assert.True(t, foundGauge, "Gauge metric not found in GetModels result")
}

// TestSetModel - тестируем пакетную вставку SetModel
func TestSetModel(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Создаем тестовые метрики
	delta := int64(100)
	value := float64(2.5)
	testModels := []models.Metrics{
		{
			ID:    "set_counter",
			MType: "counter",
			Delta: &delta,
		},
		{
			ID:    "set_gauge",
			MType: "gauge",
			Value: &value,
		},
	}

	// Вызываем тестируемую функцию
	service.SetModel(testModels)

	// Проверяем, что метрики были установлены
	counterValue, err := service.GetValue("counter", "set_counter")
	assert.NoError(t, err)
	assert.Equal(t, "100", counterValue)

	gaugeValue, err := service.GetValue("gauge", "set_gauge")
	assert.NoError(t, err)
	assert.Equal(t, "2.5", gaugeValue)
}

// TestSetModelEmptySlice - тестируем пустую пакетную вставку
func TestSetModelEmptySlice(t *testing.T) {
	memStorage := models.NewMemStorage()
	service := NewServeService(memStorage)

	// Тестируем с пустым слайсом
	service.SetModel([]models.Metrics{})

	// Проверяем, что GetAll возвращает пустую карту
	allMetrics := service.GetAll()
	assert.Empty(t, allMetrics)
}
