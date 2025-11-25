package service

import (
	"metrics/internal/agent/sender"
	models "metrics/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentService(t *testing.T) {
	metricsModel := models.NewMemStorage()

	url := "http://localhost:8080/"
	send := sender.NewSender(&url)
	rateLimit := 10
	agent := NewAgentService(metricsModel, send, rateLimit)

	require.NotNil(t, agent)

	// Проверяем, что поля структуры инициализированы
	assert.NotNil(t, agent.model)
	assert.NotNil(t, agent.rateLimit)
	assert.NotNil(t, agent.metricsCollector)
	assert.NotNil(t, agent.metricsSender)

	assert.Equal(t, agent.rateLimit, rateLimit)
	assert.Equal(t, agent.model, metricsModel)

}
