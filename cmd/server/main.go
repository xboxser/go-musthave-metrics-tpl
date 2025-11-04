package main

import (
	"metrics/internal/handler"
	models "metrics/internal/model"
	"metrics/internal/service"
	// _ "net/http/pprof" // подключаем пакет pprof
)

func main() {
	model := models.NewMemStorage()
	service := service.NewServeService(model)
	handler.Run(service)
}
