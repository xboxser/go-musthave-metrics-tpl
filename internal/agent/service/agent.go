package service

import (
	"metrics/internal/agent/sender"
	models "metrics/internal/model"
)

// generate:reset

type AgentService struct {
	metricsCollector *MetricsCollector
	metricsSender    *MetricsSender
	model            *models.MemStorage
	send             *sender.Sender
	rateLimit        int
}

func NewAgentService(model *models.MemStorage, metricsSender *MetricsSender, rateLimit int) *AgentService {
	return &AgentService{
		model:            model,
		rateLimit:        rateLimit,
		metricsCollector: NewMetricsCollector(model),
		metricsSender:    metricsSender,
	}
}

func (s *AgentService) CheckRuntime() {
	s.metricsCollector.CheckRuntime()
}

func (s *AgentService) SendMetrics() error {
	return s.metricsSender.SendMetrics(s.model)
}
