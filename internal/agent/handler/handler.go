package agent

import (
	"fmt"
	"metrics/internal/agent/config"
	"metrics/internal/agent/sender"
	"metrics/internal/agent/service"
	models "metrics/internal/model"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run() {

	configAgent := config.NewConfigAgent()

	metricsModel := models.NewMemStorage()
	send := sender.NewSender(&configAgent.URL)

	// Инициализируем gRPC если указан адрес
	if configAgent.GRPCAddress != "" {
		err := send.InitGRPC(configAgent.GRPCAddress)
		if err != nil {
			panic(fmt.Sprintf("Failed to initialize gRPC: %v", err))
		}
		defer send.CloseGRPC()
	}

	if configAgent.KEY != "" {
		send.InitHasher(configAgent.KEY)
	}

	if configAgent.CryptoKeyPath != "" {
		err := send.InitCryptoCertificate(configAgent.CryptoKeyPath)
		if err != nil {
			panic(err)
		}
	}

	metricsSender := service.NewMetricsSender(send, configAgent.RateLimit)
	// Включаем gRPC если указан адрес
	if configAgent.GRPCAddress != "" {
		metricsSender.EnableGRPC(configAgent.GRPCAddress)
	}

	service := service.NewAgentService(metricsModel, metricsSender, configAgent.RateLimit)

	pollTicker := time.NewTicker(time.Duration(configAgent.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(configAgent.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()
	// канал для получения сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-pollTicker.C:
			service.CheckRuntime()
		case <-reportTicker.C:
			_ = service.SendMetrics()
		case <-sigChan:
			_ = service.SendMetrics()
			fmt.Println("Получен сигнал завершения работы, метрики отправлены")
			return
		}
	}
}
