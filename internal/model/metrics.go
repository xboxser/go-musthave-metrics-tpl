package models

import "sync"

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
	gaugeMu sync.RWMutex
	countMu sync.RWMutex
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
	m.countMu.Lock()
	defer m.countMu.Unlock()
	m.Counter[name] += val
}

func (m *MemStorage) UpdateGauge(name string, val float64) {
	m.gaugeMu.Lock()
	defer m.gaugeMu.Unlock()
	m.Gauge[name] = val
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	m.gaugeMu.RLock()
	defer m.gaugeMu.RUnlock()
	val, ok := m.Gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	m.countMu.RLock()
	defer m.countMu.RUnlock()
	val, ok := m.Counter[name]
	return val, ok
}

func (m *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	m.gaugeMu.RLock()
	m.countMu.RLock()
	defer m.gaugeMu.RUnlock()
	defer m.countMu.RUnlock()
	return m.Gauge, m.Counter
}
