package main

import (
	"fmt"
	"metrics/internal/handler"
	models "metrics/internal/model"
	"metrics/internal/service"
)

func main() {
	fmt.Println("Hello world")

	model := models.NewMemStorage()
	service := service.NewServeService(model)
	handler.Run(service)
}
