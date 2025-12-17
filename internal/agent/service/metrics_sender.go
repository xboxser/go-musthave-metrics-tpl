package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"metrics/internal/agent/sender"
	models "metrics/internal/model"
	"metrics/internal/proto"
	"net"
	"os"
)

type MetricsSender struct {
	send      *sender.Sender
	rateLimit int
	useGRPC   bool
}

func NewMetricsSender(send *sender.Sender, rateLimit int) *MetricsSender {
	return &MetricsSender{
		send:      send,
		rateLimit: rateLimit,
	}
}

func (s *MetricsSender) EnableGRPC(grpcAddress string) error {

	err := s.send.InitGRPC(grpcAddress)
	if err != nil {
		return err
	}
	s.useGRPC = true
	return nil
}

func (s *MetricsSender) SendMetrics(storage *models.MemStorage) error {
	if s.useGRPC {
		return s.sendMetricsGRPC(storage)
	}
	return s.sendMetricsHTTP(storage)
}

func (s *MetricsSender) sendMetricsGRPC(storage *models.MemStorage) error {
	var protoMetrics []*proto.Metric

	storage.GaugeMu.RLock()
	for name, value := range storage.Gauge {
		metric := &proto.Metric{
			Id:    name,
			Type:  proto.Metric_GAUGE,
			Value: value,
		}
		protoMetrics = append(protoMetrics, metric)
	}
	storage.GaugeMu.RUnlock()

	storage.CountMu.RLock()
	for name, value := range storage.Counter {
		metric := &proto.Metric{
			Id:    name,
			Type:  proto.Metric_COUNTER,
			Delta: value,
		}
		protoMetrics = append(protoMetrics, metric)
	}
	storage.CountMu.RUnlock()

	if len(protoMetrics) == 0 {
		return nil
	}

	// Получаем IP клиента
	clientIP := s.getClientIP()

	// Отправляем метрики через gRPC
	err := s.send.SendMetricsGRPC(protoMetrics, clientIP)
	if err != nil {
		return fmt.Errorf("error sending metrics via gRPC: %w", err)
	}

	return nil
}

func (s *MetricsSender) getClientIP() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return ipv4.String()
		}
	}
	return "127.0.0.1"
}

func (s *MetricsSender) sendMetricsHTTP(storage *models.MemStorage) error {
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

	batchSize := len(metricsBatch) / s.rateLimit

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
