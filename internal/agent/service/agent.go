package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	agentModel "metrics/internal/agent/model"
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
	case uint32:
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

	chanGaugeMemMetrics := generatorGaugeMemoryMetrics()
	chanCounterMetrics := generatorCounterMetrics()

	s.fanIn([]chan agentModel.ChanCounter{chanCounterMetrics}, []chan agentModel.ChanGauge{chanGaugeMemMetrics})
}

func (s *AgentService) SendMetrics() error {
	var metrics models.Metrics
	var errs []error
	var metricsBatch []models.Metrics

	for name, value := range s.model.Gauge {
		metrics.ID = name
		metrics.MType = models.Gauge
		metrics.Value = &value
		metrics.Delta = nil
		metricsBatch = append(metricsBatch, metrics)
	}

	for name, value := range s.model.Counter {
		metrics.ID = name
		metrics.MType = models.Counter
		metrics.Value = nil
		metrics.Delta = &value
		metricsBatch = append(metricsBatch, metrics)
	}

	if len(metricsBatch) == 0 {
		return nil
	}

	resp, err := json.Marshal(metricsBatch)
	if err != nil {
		return fmt.Errorf("error json metricsBatch: %v", err)
	}
	err = s.send.SendRequest(resp)
	if err != nil {
		log.Println("error send metricsBatch", err)
		errs = append(errs, fmt.Errorf("error send metricsBatch: %v", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("failed to send metrics: %w", errors.Join(errs...))
	}
	return nil
}

func (s *AgentService) fanIn(counterChans []chan agentModel.ChanCounter, gaugeChans []chan agentModel.ChanGauge) {
	for _, counterChan := range counterChans {
		go func(c chan agentModel.ChanCounter) {
			for metric := range c {
				s.addCounter(metric.Name, metric.Value)
			}
		}(counterChan)
	}
	for _, gaugeChan := range gaugeChans {
		go func(c chan agentModel.ChanGauge) {
			for metric := range c {
				s.addGauge(metric.Name, metric.Value)
			}
		}(gaugeChan)
	}
}

func generatorCounterMetrics() chan agentModel.ChanCounter {
	outMetrics := make(chan agentModel.ChanCounter)
	go func() {
		defer close(outMetrics)
		outMetrics <- agentModel.ChanCounter{
			Name:  "PollCount",
			Value: 1,
		}
	}()
	return outMetrics
}

func generatorGaugeMemoryMetrics() chan agentModel.ChanGauge {
	outMetrics := make(chan agentModel.ChanGauge)

	go func() {
		defer close(outMetrics)
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		outMetrics <- agentModel.ChanGauge{
			Name:  "Alloc",
			Value: mem.Alloc,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "BuckHashSys",
			Value: mem.BuckHashSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "Frees",
			Value: mem.Frees,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "GCCPUFraction",
			Value: mem.GCCPUFraction,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "GCSys",
			Value: mem.GCSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapAlloc",
			Value: mem.HeapAlloc,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapIdle",
			Value: mem.HeapIdle,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapInuse",
			Value: mem.HeapInuse,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapObjects",
			Value: mem.HeapObjects,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapReleased",
			Value: mem.HeapReleased,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "HeapSys",
			Value: mem.HeapSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "LastGC",
			Value: mem.LastGC,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "Lookups",
			Value: mem.Lookups,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "MCacheInuse",
			Value: mem.MCacheInuse,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "MCacheSys",
			Value: mem.MCacheSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "MSpanInuse",
			Value: mem.MSpanInuse,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "MSpanSys",
			Value: mem.MSpanSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "Mallocs",
			Value: mem.Mallocs,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "NextGC",
			Value: mem.NextGC,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "OtherSys",
			Value: mem.OtherSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "PauseTotalNs",
			Value: mem.PauseTotalNs,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "StackInuse",
			Value: mem.StackInuse,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "StackSys",
			Value: mem.StackSys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "Sys",
			Value: mem.Sys,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "TotalAlloc",
			Value: mem.TotalAlloc,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "NumForcedGC",
			Value: mem.NumForcedGC,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "NumGC",
			Value: mem.NumGC,
		}
		outMetrics <- agentModel.ChanGauge{
			Name:  "RandomValue",
			Value: rand.Float64(),
		}
	}()
	return outMetrics
}
