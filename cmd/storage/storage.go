package storage

import (
	"fmt"
)

type MetricType string

const (
	GaugeMetric   MetricType = "gauge"
	CounterMetric MetricType = "counter"
)

type Metric struct {
	Type  MetricType
	Name  string
	Value interface{}
}

type MemStorage struct {
	data map[string]Metric
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		data: make(map[string]Metric),
	}
}

func (s *MemStorage) UpdateMerics(m Metric) error {

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
	metric, exist := s.data[name]
	return metric, exist
}
