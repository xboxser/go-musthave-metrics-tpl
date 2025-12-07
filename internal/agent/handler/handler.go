package agent

import (
	"metrics/internal/agent/config"
	"metrics/internal/agent/sender"
	"metrics/internal/agent/service"
	models "metrics/internal/model"
	"time"
)

func Run() {

	configAgent := config.NewConfigAgent()

	metricsModel := models.NewMemStorage()
	send := sender.NewSender(&configAgent.URL)
	if configAgent.KEY != "" {
		send.InitHasher(configAgent.KEY)
	}

	if configAgent.CryptoKeyPath != "" {
		err := send.InitCryptoCertificate(configAgent.CryptoKeyPath)
		if err != nil {
			panic(err)
		}
	}
	service := service.NewAgentService(metricsModel, send, configAgent.RateLimit)

	pollTicker := time.NewTicker(time.Duration(configAgent.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(configAgent.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			service.CheckRuntime()
		case <-reportTicker.C:
			_ = service.SendMetrics()
		}
	}
}
