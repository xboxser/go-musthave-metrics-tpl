package service

import (
	"math/rand/v2"
	agentModel "metrics/internal/agent/model"
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func sendMemStats(outMetrics chan agentModel.ChanGauge, mem *runtime.MemStats) {
	outMetrics <- agentModel.ChanGauge{
		Name:  "Alloc",
		Value: mem.Alloc,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "BuckHashSys",
		Value: mem.BuckHashSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "Frees",
		Value: mem.Frees,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "GCCPUFraction",
		Value: mem.GCCPUFraction,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "GCSys",
		Value: mem.GCSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapAlloc",
		Value: mem.HeapAlloc,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapIdle",
		Value: mem.HeapIdle,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapInuse",
		Value: mem.HeapInuse,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapObjects",
		Value: mem.HeapObjects,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapReleased",
		Value: mem.HeapReleased,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "HeapSys",
		Value: mem.HeapSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "LastGC",
		Value: mem.LastGC,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "Lookups",
		Value: mem.Lookups,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "MCacheInuse",
		Value: mem.MCacheInuse,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "MCacheSys",
		Value: mem.MCacheSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "MSpanInuse",
		Value: mem.MSpanInuse,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "MSpanSys",
		Value: mem.MSpanSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "Mallocs",
		Value: mem.Mallocs,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "NextGC",
		Value: mem.NextGC,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "OtherSys",
		Value: mem.OtherSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "PauseTotalNs",
		Value: mem.PauseTotalNs,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "StackInuse",
		Value: mem.StackInuse,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "StackSys",
		Value: mem.StackSys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "Sys",
		Value: mem.Sys,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "TotalAlloc",
		Value: mem.TotalAlloc,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "NumForcedGC",
		Value: mem.NumForcedGC,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "NumGC",
		Value: mem.NumGC,
	}
}

func sendRandomValue(outMetrics chan agentModel.ChanGauge) {
	outMetrics <- agentModel.ChanGauge{
		Name:  "RandomValue",
		Value: rand.Float64(),
	}
}

func sendGaugeGopsutil(outMetrics chan agentModel.ChanGauge) error {
	mem, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "TotalMemory",
		Value: mem.Total,
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "FreeMemory",
		Value: mem.Free,
	}
	cpuPercent, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}
	for i, percent := range cpuPercent {
		outMetrics <- agentModel.ChanGauge{
			Name:  "CPUutilization" + strconv.Itoa(i),
			Value: percent,
		}
	}
	outMetrics <- agentModel.ChanGauge{
		Name:  "RandomValue",
		Value: rand.Float64(),
	}
	return nil
}
