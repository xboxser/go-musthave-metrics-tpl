package service

import (
	models "metrics/internal/model"
	"testing"
)

// _ "net/http/pprof" // подключаем пакет pprof

func TestNewServeService(t *testing.T) {
	model := models.NewMemStorage()
	service := NewServeService(model)

	if service == nil {
		t.Fatal("failed to create service")
	}
	if service.model != model {
		t.Fatal("NewSender returned wrong model")
	}
}
