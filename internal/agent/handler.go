package agent

import (
	"fmt"
	"math/rand/v2"
	models "metrics/internal/model"
	"metrics/internal/service"
	"runtime"
	"time"
)

func Run() {
	metricsModel := models.NewMemStorage()
	service := service.NewAgentService(metricsModel)

	i := 0
	for {
		readRuntime(service)
		i += 2
		time.Sleep(2 * time.Second)

		if i%10 == 0 {
			fmt.Println("send request server", i)
		}

		if i > 20 {
			break
		}
	}
	service.Print()
}

func readRuntime(service *service.AgentService) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	service.AddGauge("Alloc", mem.Alloc)
	service.AddGauge("BuckHashSys", mem.BuckHashSys)
	service.AddGauge("Frees", mem.Frees)
	service.AddGauge("GCCPUFraction", mem.GCCPUFraction)
	service.AddGauge("GCSys", mem.GCSys)
	service.AddGauge("HeapAlloc", mem.HeapAlloc)
	service.AddGauge("HeapIdle", mem.HeapIdle)
	service.AddGauge("HeapInuse", mem.HeapInuse)
	service.AddGauge("HeapObjects", mem.HeapObjects)
	service.AddGauge("HeapReleased", mem.HeapReleased)
	service.AddGauge("HeapSys", mem.HeapSys)
	service.AddGauge("LastGC", mem.LastGC)
	service.AddGauge("Lookups", mem.Lookups)
	service.AddGauge("MCacheInuse", mem.MCacheInuse)
	service.AddGauge("MCacheSys", mem.MCacheSys)
	service.AddGauge("MSpanInuse", mem.MSpanInuse)
	service.AddGauge("MSpanSys", mem.MSpanSys)
	service.AddGauge("Mallocs", mem.Mallocs)
	service.AddGauge("NextGC", mem.NextGC)
	service.AddGauge("OtherSys", mem.OtherSys)
	service.AddGauge("PauseTotalNs", mem.PauseTotalNs)
	service.AddGauge("StackInuse", mem.StackInuse)
	service.AddGauge("StackSys", mem.StackSys)
	service.AddGauge("Sys", mem.Sys)
	service.AddGauge("TotalAlloc", mem.TotalAlloc)

	service.AddGauge("RandomValue", rand.Float64())
	service.AddCounter("PollCount", 1)
}
