package models

import "sync"

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	GaugeMu sync.RWMutex
	CountMu sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

type Storage interface {
	UpdateGauge(name string, val float64)
	UpdateCounter(name string, val int64)
	GetGauge(name string) (float64, bool)
	GetCounter(name string) (int64, bool)
	GetAll() (map[string]float64, map[string]int64)
}

func (m *MemStorage) UpdateCounter(name string, val int64) {
	m.CountMu.Lock()
	defer m.CountMu.Unlock()
	m.Counter[name] += val
}

func (m *MemStorage) UpdateGauge(name string, val float64) {
	m.GaugeMu.Lock()
	defer m.GaugeMu.Unlock()
	m.Gauge[name] = val
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.GaugeMu.RLock()
	defer m.GaugeMu.RUnlock()
	val, ok := m.Gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.CountMu.RLock()
	defer m.CountMu.RUnlock()
	val, ok := m.Counter[name]
	return val, ok
}

func (m *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	m.GaugeMu.RLock()
	m.CountMu.RLock()

	// Создаем копии map'ов для безопасного возврата
	gaugeCopy := make(map[string]float64)
	for k, v := range m.Gauge {
		gaugeCopy[k] = v
	}

	counterCopy := make(map[string]int64)
	for k, v := range m.Counter {
		counterCopy[k] = v
	}

	defer m.CountMu.RUnlock()
	defer m.GaugeMu.RUnlock()
	return gaugeCopy, counterCopy
}
