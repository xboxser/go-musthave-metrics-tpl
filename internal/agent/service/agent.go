package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand/v2"
	"metrics/internal/agent/sender"
	models "metrics/internal/model"
	"runtime"
)

type AgentService struct {
	model *models.MemStorage
	send  *sender.Sender
}

func NewAgentService(model *models.MemStorage, send *sender.Sender) *AgentService {
	return &AgentService{
		model: model,
		send:  send,
	}
}

func (s *AgentService) addGauge(name string, val interface{}) error {
	var floatVal float64

	switch v := val.(type) {
	case float64:
		floatVal = v
	case int:
		floatVal = float64(v)
	case int64:
		floatVal = float64(v)
	case uint64:
		floatVal = float64(v)
	case float32:
		floatVal = float64(v)
	case bool:
		if v {
			floatVal = 1.
		} else {
			floatVal = 0.
		}
	default:
		return errors.New("error update operation: incorrect type value")
	}
	s.model.Gauge[name] = floatVal
	return nil
}

func (s *AgentService) addCounter(name string, val int64) error {
	s.model.Counter[name] += val
	return nil
}

func (s *AgentService) Print() {
	fmt.Println(s.model)
}

func (s *AgentService) CheckRuntime() {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	s.addGauge("Alloc", mem.Alloc)
	s.addGauge("BuckHashSys", mem.BuckHashSys)
	s.addGauge("Frees", mem.Frees)
	s.addGauge("GCCPUFraction", mem.GCCPUFraction)
	s.addGauge("GCSys", mem.GCSys)
	s.addGauge("HeapAlloc", mem.HeapAlloc)
	s.addGauge("HeapIdle", mem.HeapIdle)
	s.addGauge("HeapInuse", mem.HeapInuse)
	s.addGauge("HeapObjects", mem.HeapObjects)
	s.addGauge("HeapReleased", mem.HeapReleased)
	s.addGauge("HeapSys", mem.HeapSys)
	s.addGauge("LastGC", mem.LastGC)
	s.addGauge("Lookups", mem.Lookups)
	s.addGauge("MCacheInuse", mem.MCacheInuse)
	s.addGauge("MCacheSys", mem.MCacheSys)
	s.addGauge("MSpanInuse", mem.MSpanInuse)
	s.addGauge("MSpanSys", mem.MSpanSys)
	s.addGauge("Mallocs", mem.Mallocs)
	s.addGauge("NextGC", mem.NextGC)
	s.addGauge("OtherSys", mem.OtherSys)
	s.addGauge("PauseTotalNs", mem.PauseTotalNs)
	s.addGauge("StackInuse", mem.StackInuse)
	s.addGauge("StackSys", mem.StackSys)
	s.addGauge("Sys", mem.Sys)
	s.addGauge("TotalAlloc", mem.TotalAlloc)
	s.addGauge("NumForcedGC", mem.NumForcedGC)
	s.addGauge("NumForcedGC", mem.NumGC)

	s.addGauge("RandomValue", rand.Float64())
	s.addCounter("PollCount", 1)
}

func (s *AgentService) SendMetrics() error {
	var metrics models.Metrics
	var errs []error

	for name, value := range s.model.Gauge {

		metrics.ID = name
		metrics.MType = models.Gauge
		metrics.Value = &value
		metrics.Delta = nil
		resp, err := json.Marshal(metrics)
		if err != nil {
			errs = append(errs, fmt.Errorf("error json Gauge: %v", err))
			continue
		}
		err = s.send.SendRequest(resp)
		if err != nil {
			errs = append(errs, fmt.Errorf("error send Gauge: %v", err))
		}
	}

	for name, value := range s.model.Counter {
		metrics.ID = name
		metrics.MType = models.Counter
		metrics.Value = nil
		metrics.Delta = &value
		resp, err := json.Marshal(metrics)
		if err != nil {
			errs = append(errs, fmt.Errorf("error json Counter: %v", err))
			continue
		}
		err = s.send.SendRequest(resp)
		if err != nil {
			errs = append(errs, fmt.Errorf("error send Counter: %v", err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send metrics: %w", errors.Join(errs...))
	}
	return nil
}
