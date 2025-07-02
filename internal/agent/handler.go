package agent

import (
	"flag"
	"metrics/internal/agent/sender"
	"metrics/internal/agent/service"
	models "metrics/internal/model"
	"os"
	"time"
)

func Run() {

	agentFlags := flag.NewFlagSet("agent", flag.ExitOnError)

	reportInterval := agentFlags.Int("r", 10, "The interval for sending data to the server")
	pollInterval := agentFlags.Int("p", 2, "The interval for building metrics")
	url := agentFlags.String("a", "localhost:8080", "port server")
	agentFlags.Parse(os.Args[1:])

	metricsModel := models.NewMemStorage()
	send := sender.NewSender(url)
	service := service.NewAgentService(metricsModel, send)

	i := 0
	for {
		if i%*pollInterval == 0 {
			service.CheckRuntime()
		}

		i += 1
		time.Sleep(1 * time.Second)

		if i%*reportInterval == 0 {
			service.SendMetrics()
		}
	}
}
