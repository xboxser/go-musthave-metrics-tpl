package agent

import (
	"flag"
	"fmt"
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

	pollTicker := time.NewTicker(time.Duration(*pollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(*reportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			service.CheckRuntime()
		case <-reportTicker.C:
			err := service.SendMetrics()
			if err != nil {
				fmt.Printf("fail to send %v", err)
				return
			}
		}
	}
}
