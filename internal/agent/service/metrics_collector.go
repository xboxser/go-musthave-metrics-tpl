package service

import (
	"errors"
	"fmt"
	agentModel "metrics/internal/agent/model"
	models "metrics/internal/model"
)

type MetricsCollector struct {
	model *models.MemStorage
}

func NewMetricsCollector(model *models.MemStorage) *MetricsCollector {
	return &MetricsCollector{model: model}
}

func (s *MetricsCollector) CheckRuntime() {

	chanGaugeMemMetrics := generatorGaugeMemoryMetrics()
	chanCounterMetrics := generatorCounterMetrics()
	chanGaugeGopMetrics := generatorGaugeGopsutilMetrics()

	s.fanIn([]chan agentModel.ChanCounter{chanCounterMetrics}, []chan agentModel.ChanGauge{chanGaugeMemMetrics, chanGaugeGopMetrics})
}

func (s *MetricsCollector) addGauge(name string, val interface{}) error {
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
	s.model.UpdateGauge(name, floatVal)
	return nil
}

func (s *MetricsCollector) addCounter(name string, val int64) error {
	s.model.UpdateCounter(name, val)
	return nil
}

func (s *MetricsCollector) Print() {
	fmt.Println(s.model)
}

func (s *MetricsCollector) fanIn(counterChans []chan agentModel.ChanCounter, gaugeChans []chan agentModel.ChanGauge) {
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
