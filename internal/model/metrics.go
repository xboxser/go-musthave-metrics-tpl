package models

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

// NOTE: Не усложняем пример, вводя иерархическую вложенность структур.
// Органичиваясь плоской моделью.
// Delta и Value объявлены через указатели,
// что бы отличать значение "0", от не заданного значения
// и соответственно не кодировать в структуру.
type MetricsDefault struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	Hash  string   `json:"hash,omitempty"`
}

type MemStorage struct {
	Gauge   map[string]float64
	Counter map[string]int64
}

func NewMemStorage() MemStorage {
	return MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}
}

func (s MemStorage) Update(t string, name string, val string) error {
	if t == Counter {
		val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.Counter[name] += val
	} else if t == Gauge {
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.Gauge[name] = val
	} else {
		return errors.New("error update operation: incorrect type")
	}
	fmt.Println(s)
	return nil
}
