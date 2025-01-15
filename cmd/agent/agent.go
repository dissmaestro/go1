package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

type Metric struct {
	Gauge     map[string]float64
	Counter   map[string]int64
	PollCount int64
}

func collectRuntimeMetrics() Metric {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	metrics := Metric{
		Gauge: map[string]float64{
			"Alloc":         float64(m.Alloc),
			"BuckHashSys":   float64(m.BuckHashSys),
			"Frees":         float64(m.Frees),
			"GCCPUFraction": m.GCCPUFraction,
			"GCSys":         float64(m.GCSys),
			"HeapAlloc":     float64(m.HeapAlloc),
			"HeapIdle":      float64(m.HeapIdle),
			"HeapInuse":     float64(m.HeapInuse),
			"HeapObjects":   float64(m.HeapObjects),
			"HeapReleased":  float64(m.HeapReleased),
			"HeapSys":       float64(m.HeapSys),
			"LastGC":        float64(m.LastGC),
			"Lookups":       float64(m.Lookups),
			"MCacheInuse":   float64(m.MCacheInuse),
			"MCacheSys":     float64(m.MCacheSys),
			"MSpanInuse":    float64(m.MSpanInuse),
			"MSpanSys":      float64(m.MSpanSys),
			"Mallocs":       float64(m.Mallocs),
			"NextGC":        float64(m.NextGC),
			"NumForcedGC":   float64(m.NumForcedGC),
			"NumGC":         float64(m.NumGC),
			"OtherSys":      float64(m.OtherSys),
			"PauseTotalNs":  float64(m.PauseTotalNs),
			"StackInuse":    float64(m.StackInuse),
			"StackSys":      float64(m.StackSys),
			"Sys":           float64(m.Sys),
			"TotalAlloc":    float64(m.TotalAlloc),
			"RandomValue":   rand.Float64() * 100, // Случайное значение для демонстрации
		},
		Counter:   map[string]int64{},
		PollCount: atomic.AddInt64(new(int64), 1),
	}

	return metrics
}

func RunAgent() {
	var pollInterval time.Duration = 2 * time.Second
	var reportInterval time.Duration = 10 * time.Second
	var serverUrl string = "http://localhost:8080/update"

	metricsChan := make(chan Metric, 1)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			metrics := collectRuntimeMetrics()
			metricsChan <- metrics
			fmt.Println("Metrics collected:", metrics.Gauge)
		}
	}()

	go func() {
		for range time.Tick(reportInterval) {
			metrics := <-metricsChan
			for name, value := range metrics.Gauge {
				url := fmt.Sprintf("%s/gauge/%s%d", serverUrl, name, value)
				_, err := http.Post(url, "text/plain", nil)
				if err != nil {
					log.Println("Failed to send gauge metric:", err)
				}
			}
			url := fmt.Sprintf("%s/counter/PollCount/%d", serverUrl, metrics.PollCount)
			_, err := http.Post(url, "text/plain", nil)
			if err != nil {
				log.Println("Failed to send counter metric:", err)
			}
		}

	}()
}

func main() {
	RunAgent()
	select {} // Блокируем main горутину для работы агента
}
