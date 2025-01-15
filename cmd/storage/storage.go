package storage

import (
	"fmt"
	"sync"
)

const (
	GaugeMetric   string = "gauge"
	CounterMetric string = "counter"
)

type Metric struct {
	Type  string
	Name  string
	Value interface{}
}

type MemStorage struct {
	mu   sync.Mutex
	data map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]Metric),
	}
}

func (s *MemStorage) UpdateMerics(m Metric) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m.Type != CounterMetric && m.Type != GaugeMetric {
		return fmt.Errorf("invalid metric type: %s. Expected 'gauge' or 'counter'", m.Type)
	}

	if exist, ok := s.data[m.Name]; ok {
		if exist.Type != m.Type {
			return fmt.Errorf("incorrect metric type for %s", m.Name)
		}
		if m.Type == CounterMetric {
			switch v := m.Value.(type) {
			case int64:
				exist.Value = exist.Value.(int64) + v
				s.data[m.Name] = exist
			default:
				return fmt.Errorf("incorrect value for counter")
			}
		} else {
			s.data[m.Name] = m
		}
	} else {
		s.data[m.Name] = m
	}
	return nil
}

func (s *MemStorage) GetMetric(name string) (Metric, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	metric, exist := s.data[name]
	return metric, exist
}
