package main

import (
	"metrics/internal/handler"
	models "metrics/internal/model"
	"metrics/internal/service"
	// _ "net/http/pprof" // подключаем пакет pprof
)

// go test -coverprofile=coverage.out ./...
// go tool cover -func=coverage.out | grep "total:"

func main() {
	model := models.NewMemStorage()
	service := service.NewServeService(model)
	handler.Run(service)
}
