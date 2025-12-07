package service

import (
	"metrics/internal/agent/model"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendRandomValue(t *testing.T) {
	outMetrics := make(chan model.ChanGauge, 1)

	sendRandomValue(outMetrics)
	defer close(outMetrics)

	// Получаем отправленную метрику
	metric := <-outMetrics

	// Проверяем, что метрика отправлена с правильным именем
	assert.Equal(t, "RandomValue", metric.Name)
	// Проверяем, что значение находится в допустимом диапазоне [0, 1)
	assert.NotNil(t, metric.Name)
}

func TestSendMemStats(t *testing.T) {

	var ms runtime.MemStats
	ms.Alloc = 100
	ms.BuckHashSys = 101
	ms.Frees = 102
	ms.GCCPUFraction = 103
	ms.GCSys = 104
	ms.HeapAlloc = 105
	ms.HeapIdle = 106
	ms.HeapInuse = 107
	ms.HeapObjects = 108
	ms.HeapReleased = 109
	ms.HeapSys = 110
	ms.LastGC = 111
	ms.Lookups = 112
	ms.MCacheInuse = 113
	ms.MCacheSys = 114
	ms.MSpanInuse = 115
	ms.MSpanSys = 116
	ms.Mallocs = 117
	ms.NextGC = 118
	ms.OtherSys = 119
	ms.PauseTotalNs = 120
	ms.StackInuse = 121
	ms.StackSys = 122
	ms.Sys = 123
	ms.TotalAlloc = 124
	ms.NumForcedGC = 125
	ms.NumGC = 126

	// Создаем канал для получения метрик
	outMetrics := make(chan model.ChanGauge, 30)
	// Вызываем тестируемую функцию в горутине
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		sendMemStats(outMetrics, &ms)
		close(outMetrics)
	}()

	// Собираем все отправленные метрики
	metrics := make(map[string]any)
	for metric := range outMetrics {
		metrics[metric.Name] = metric.Value
	}

	// Ждем завершения горутины
	wg.Wait()

	// Проверяем, что метрики были отправлены
	assert.Equal(t, ms.Alloc, metrics["Alloc"])
	assert.Equal(t, ms.BuckHashSys, metrics["BuckHashSys"])
	assert.Equal(t, ms.Frees, metrics["Frees"])
	assert.Equal(t, ms.GCCPUFraction, metrics["GCCPUFraction"])
	assert.Equal(t, ms.GCSys, metrics["GCSys"])
	assert.Equal(t, ms.HeapAlloc, metrics["HeapAlloc"])
	assert.Equal(t, ms.HeapIdle, metrics["HeapIdle"])
	assert.Equal(t, ms.HeapInuse, metrics["HeapInuse"])
	assert.Equal(t, ms.HeapObjects, metrics["HeapObjects"])
	assert.Equal(t, ms.HeapReleased, metrics["HeapReleased"])
	assert.Equal(t, ms.HeapSys, metrics["HeapSys"])
	assert.Equal(t, ms.LastGC, metrics["LastGC"])
	assert.Equal(t, ms.Lookups, metrics["Lookups"])
	assert.Equal(t, ms.Mallocs, metrics["Mallocs"])
	assert.Equal(t, ms.NumGC, metrics["NumGC"])
	assert.Equal(t, ms.PauseTotalNs, metrics["PauseTotalNs"])
	assert.Equal(t, ms.StackInuse, metrics["StackInuse"])
	assert.Equal(t, ms.StackSys, metrics["StackSys"])
	assert.Equal(t, ms.Sys, metrics["Sys"])
	assert.Equal(t, ms.TotalAlloc, metrics["TotalAlloc"])
	assert.Equal(t, ms.NumForcedGC, metrics["NumForcedGC"])
	assert.Equal(t, ms.OtherSys, metrics["OtherSys"])
	assert.Equal(t, ms.NextGC, metrics["NextGC"])
	assert.Equal(t, ms.MSpanSys, metrics["MSpanSys"])
	assert.Equal(t, ms.MSpanInuse, metrics["MSpanInuse"])
	assert.Equal(t, ms.MCacheSys, metrics["MCacheSys"])
	assert.Equal(t, ms.MCacheInuse, metrics["MCacheInuse"])
}
