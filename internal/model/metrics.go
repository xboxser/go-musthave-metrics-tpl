package models

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
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
	m.Counter[name] += val
}

func (m *MemStorage) UpdateGauge(name string, val float64) {
	m.Gauge[name] = val
}

func (m *MemStorage) GetGauge(name string) (float64, bool) {
	val, ok := m.Gauge[name]
	return val, ok
}

func (m *MemStorage) GetCounter(name string) (int64, bool) {
	val, ok := m.Counter[name]
	return val, ok
}

func (m *MemStorage) GetAll() (map[string]float64, map[string]int64) {
	return m.Gauge, m.Counter
}
