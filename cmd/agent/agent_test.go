package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func TestCollectRuntimeMetrics(T *testing.T) {
	metrics := collectRuntimeMetrics()

	if len(metrics.Gauge) == 0 {
		T.Errorf("Expected, that metrics guage to be non-empty, got %d", len(metrics.Gauge))
	}

	if _, ok := metrics.Gauge["Alloc"]; !ok {
		T.Errorf("Expected key 'Alloc' in Gauge, but it was missing")
	}

	if metrics.Gauge["RandomValue"] < 0 || metrics.Gauge["RandomValue"] > 100 {
		T.Errorf("Expected RandomValue to be between 0 and 100, got %f", metrics.Gauge["RandomValue"])
	}

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	if metrics.Gauge["LastGC"] != float64(m.LastGC) {
		T.Errorf("Expected LastGC to be %v, got %v", m.LastGC, metrics.Gauge["LastGC"])
	}
}

func TestSendMerics(T *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			T.Errorf("Expected POST method, got %s", r.Method)
		}
		fmt.Fprintln(w, "OK")
	}))
	defer server.Close()

	metrics := collectRuntimeMetrics()
	for name, value := range metrics.Gauge {
		url := fmt.Sprintf("%s/gauge/%s/%v", server.URL, name, value)
		resp, err := http.Post(url, "text/plain", nil)
		if err != nil {
			T.Errorf("Failed to send metric: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			T.Errorf("Expected status code 200, got %d", resp.StatusCode)
		}
	}
}

func TestMetricsCollectionLoop(t *testing.T) {
	metricsChan := make(chan Metric, 1)

	go func() {
		for range time.Tick(10 * time.Millisecond) {
			metrics := collectRuntimeMetrics()
			metricsChan <- metrics
		}
	}()

	select {
	case metrics := <-metricsChan:
		if len(metrics.Gauge) == 0 {
			t.Errorf("Expected metrics.Gauge to have values, got empty map")
		}
	case <-time.After(50 * time.Millisecond):
		t.Error("Timed out waiting for metrics collection")
	}
}
