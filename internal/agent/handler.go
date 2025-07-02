package agent

import (
	models "metrics/internal/model"
	"metrics/internal/service"
	"time"
)

func Run() {
	metricsModel := models.NewMemStorage()
	service := service.NewAgentService(metricsModel)

	i := 0
	for {
		service.CheckRuntime()
		i += 2
		time.Sleep(2 * time.Second)

		if i%10 == 0 {
			service.Send()
		}
	}
}
