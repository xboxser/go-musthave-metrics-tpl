package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	agentModel "metrics/internal/agent/model"
	"metrics/internal/agent/sender"
	models "metrics/internal/model"
	"runtime"
)

type MetricsSender struct {
	send      *sender.Sender
	rateLimit int
}

func NewMetricsSender(send *sender.Sender, rateLimit int) *MetricsSender {
	return &MetricsSender{
		send:      send,
		rateLimit: rateLimit,
	}
}

func (s *MetricsSender) SendMetrics(storage *models.MemStorage) error {
	var metrics models.Metrics
	var errs []error
	var metricsBatch []models.Metrics

	storage.GaugeMu.RLock()
	for name, value := range storage.Gauge {
		metrics.ID = name
		metrics.MType = models.Gauge
		metrics.Value = &value
		metrics.Delta = nil
		metricsBatch = append(metricsBatch, metrics)
	}
	storage.GaugeMu.RUnlock()

	storage.CountMu.RLock()
	for name, value := range storage.Counter {
		metrics.ID = name
		metrics.MType = models.Counter
		metrics.Value = nil
		metrics.Delta = &value
		metricsBatch = append(metricsBatch, metrics)
	}
	storage.CountMu.RUnlock()

	if len(metricsBatch) == 0 {
		return nil
	}
	const batchSize = 5

	byteChan := make(chan []byte, len(metricsBatch))
	errorChan := make(chan error, s.rateLimit)

	for i := 0; i < s.rateLimit; i++ {
		go func(byteChan <-chan []byte, errorChan chan<- error) {
			for resp := range byteChan {
				err := s.send.SendRequest(resp)
				errorChan <- err
			}

		}(byteChan, errorChan)
	}

	countBatch := 0
	for i := 0; i < len(metricsBatch); i += batchSize {
		end := i + batchSize
		if end > len(metricsBatch) {
			end = len(metricsBatch)
		}
		resp, err := json.Marshal(metricsBatch[i:end])
		if err != nil {
			close(byteChan)
			close(errorChan)
			return fmt.Errorf("error json metricsBatch: %v", err)
		}
		countBatch++
		byteChan <- resp
	}

	close(byteChan)
	completedWorkers := 0
	for completedWorkers < countBatch {
		if err := <-errorChan; err != nil {
			errs = append(errs, fmt.Errorf("error send metricsBatch: %v", err))
		}
		completedWorkers++
	}
	close(errorChan)

	if len(errs) > 0 {
		return fmt.Errorf("failed to send metrics: %w", errors.Join(errs...))
	}
	return nil
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

func generatorGaugeGopsutilMetrics() chan agentModel.ChanGauge {
	outMetrics := make(chan agentModel.ChanGauge)
	go func() {
		defer close(outMetrics)
		err := sendGaugeGopsutil(outMetrics)
		if err != nil {
			log.Printf("Ошибка при метрик: %v", err)
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
		sendMemStats(outMetrics, &mem)
		sendRandomValue(outMetrics)
	}()
	return outMetrics
}
