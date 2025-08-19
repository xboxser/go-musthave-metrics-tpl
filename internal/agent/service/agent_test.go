package service

import (
	models "metrics/internal/model"
	"testing"
)

func TestAddGauge(t *testing.T) {
	type args struct {
		name string
		val  interface{}
	}

	tests := []struct {
		name     string
		args     args
		valError bool
		want     float64
	}{
		{name: "valid parameters", args: args{name: "name", val: float64(123)}, valError: false, want: float64(123)},
		{name: "valid parameters bool", args: args{name: "name", val: bool(true)}, valError: false, want: float64(1)},
		{name: "valid parameters", args: args{name: "name", val: "lololol"}, valError: true, want: float64(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &models.MemStorage{
				Gauge: make(map[string]float64),
			}
			s := NewMetricsCollector(storage)
			err := s.addGauge(tt.args.name, tt.args.val)

			if (err != nil) != tt.valError {
				t.Errorf("addGauge error %v", err)
				return
			}

			if err == nil {
				if val, ok := storage.GetGauge(tt.args.name); !ok {
					t.Errorf("addGauge value not found in storage")
				} else if val != tt.want {
					t.Errorf("addGauge = %v, want %v", val, tt.want)
				}
			}
		})
	}
}

func TestAddCounter(t *testing.T) {
	type args struct {
		name string
		val  int64
	}

	tests := []struct {
		name     string
		args     args
		valError bool
		want     int64
	}{
		{name: "valid parameters", args: args{name: "name", val: int64(123)}, valError: false, want: int64(123)},
		{name: "valid parameters bool", args: args{name: "name", val: int64(1)}, valError: false, want: int64(1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &models.MemStorage{
				Counter: make(map[string]int64),
			}
			s := NewMetricsCollector(storage)
			err := s.addCounter(tt.args.name, tt.args.val)

			if (err != nil) != tt.valError {
				t.Errorf("addCounter error: %v", err)
				return
			}

			if err == nil {
				if val, ok := storage.GetCounter(tt.args.name); !ok {
					t.Errorf("addCounter value not found in storage")
				} else if val != tt.want {
					t.Errorf("addCounter = %v, want %v", val, tt.want)
				}
			}
		})
	}
}
