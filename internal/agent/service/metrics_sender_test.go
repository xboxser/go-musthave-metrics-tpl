package service

import (
	"metrics/internal/agent/sender"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsSender(t *testing.T) {

	url := "http://localhost:8080/"
	send := sender.NewSender(&url)
	rateLimit := 10
	ms := NewMetricsSender(send, rateLimit)

	require.NotNil(t, ms)

	// Проверяем, что поля структуры инициализированы
	assert.NotNil(t, ms.send)
	assert.NotNil(t, ms.rateLimit)

	assert.Equal(t, ms.send, send)
	assert.Equal(t, ms.rateLimit, rateLimit)
}
