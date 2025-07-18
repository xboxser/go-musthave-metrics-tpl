package agent

import (
	"metrics/internal/agent/sender"
	"metrics/internal/agent/service"
	models "metrics/internal/model"
	"time"
)

func Run() {

	configAgent := newCongigAgent()

	metricsModel := models.NewMemStorage()
	send := sender.NewSender(&configAgent.URL)
	service := service.NewAgentService(metricsModel, send)

	pollTicker := time.NewTicker(time.Duration(configAgent.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(configAgent.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			service.CheckRuntime()
		case <-reportTicker.C:
			err := service.SendMetrics()
			if err != nil {
				// fmt.Printf("fail to send %v", err)
				// return
			}
		}
	}
}
