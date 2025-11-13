package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateCounter(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  int64
	}{
		{name: "simple number", key: "test_key", val: 10},
		{name: "negative number", key: "test_key", val: -10},
		{name: "zero", key: "test_key", val: 0},
		{name: "big number", key: "test_key", val: 52315611555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewMemStorage()

			metrics.UpdateCounter(tt.key, tt.val)
			val := metrics.Counter[tt.key]
			if val != tt.val {
				t.Errorf("expected %d, got %d", tt.val, val)
			}

			metrics.UpdateCounter(tt.key, tt.val)
			val = metrics.Counter[tt.key]
			if val != tt.val+tt.val {
				t.Errorf("double expected %d, got %d", tt.val+tt.val, val)
			}
		})
	}
}

func TestUpdateGauge(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  float64
	}{
		{name: "simple number", key: "test_key", val: 10},
		{name: "negative number", key: "test_key", val: -10},
		{name: "zero", key: "test_key", val: 0},
		{name: "big number", key: "test_key", val: 52315611555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewMemStorage()

			metrics.UpdateGauge(tt.key, tt.val)
			val := metrics.Gauge[tt.key]
			if val != tt.val {
				t.Errorf("expected %f, got %f", tt.val, val)
			}

			metrics.UpdateGauge(tt.key, tt.val)
			val = metrics.Gauge[tt.key]
			if val != tt.val {
				t.Errorf("double expected %f, got %f", tt.val, val)
			}
		})
	}
}

func TestGetGauge(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  float64
	}{
		{name: "simple number", key: "test_key", val: 10},
		{name: "negative number", key: "test_key", val: -10},
		{name: "zero", key: "test_key", val: 0},
		{name: "big number", key: "test_key", val: 52315611555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewMemStorage()

			val, ok := metrics.GetGauge(tt.key)

			assert.Equal(t, ok, false)
			assert.Equal(t, val, 0.)

			metrics.UpdateGauge(tt.key, tt.val)
			val, ok = metrics.GetGauge(tt.key)
			if val != tt.val || ok == false {
				t.Errorf("expected %f, got %f", tt.val, val)
			}

			metrics.UpdateGauge(tt.key, tt.val)
			val, ok = metrics.GetGauge(tt.key)
			if val != tt.val || ok == false {
				t.Errorf("double expected %f, got %f", tt.val, val)
			}
		})
	}
}

func TestGetCounter(t *testing.T) {
	tests := []struct {
		name string
		key  string
		val  int64
	}{
		{name: "simple number", key: "test_key", val: 10},
		{name: "negative number", key: "test_key", val: -10},
		{name: "zero", key: "test_key", val: 0},
		{name: "big number", key: "test_key", val: 52315611555},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewMemStorage()

			val, ok := metrics.GetCounter(tt.key)

			assert.Equal(t, ok, false)
			assert.Equal(t, val, int64(0))
			metrics.UpdateCounter(tt.key, tt.val)

			val, ok = metrics.GetCounter(tt.key)
			if val != tt.val || ok == false {
				t.Errorf("expected %d, got %d", tt.val, val)
			}

			metrics.UpdateCounter(tt.key, tt.val)
			val, ok = metrics.GetCounter(tt.key)
			if val != tt.val+tt.val || ok == false {
				t.Errorf("double expected %d, got %d", tt.val+tt.val, val)
			}
		})
	}
}

func TestGetAll(t *testing.T) {

	type counter struct {
		key string
		val int64
	}

	type gauge struct {
		key string
		val float64
	}
	tests := []struct {
		name    string
		counter []counter
		gauge   []gauge
	}{
		{
			name: "simple number",
			counter: []counter{
				counter{key: "test_key", val: 10},
				counter{key: "empty", val: 10},
			},
			gauge: []gauge{
				gauge{key: "test_key", val: 10},
				gauge{key: "empty", val: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := NewMemStorage()

			for _, c := range tt.counter {
				metrics.UpdateCounter(c.key, c.val)
			}
			for _, g := range tt.gauge {
				metrics.UpdateGauge(g.key, g.val)
			}

			gauge, counter := metrics.GetAll()
			for _, c := range tt.counter {
				if counter[c.key] != c.val {
					t.Errorf("expected %d, got %d", c.val, counter[c.key])
				}
			}
			for _, g := range tt.gauge {
				if gauge[g.key] != g.val {
					t.Errorf("expected %f, got %f", g.val, gauge[g.key])
				}
			}
		})
	}
}
