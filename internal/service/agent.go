package service

import (
	"errors"
	"fmt"
	models "metrics/internal/model"
)

type AgentService struct {
	model *models.MemStorage
}

func NewAgentService(model *models.MemStorage) *AgentService {
	return &AgentService{
		model: model,
	}
}

func (s *AgentService) AddGauge(name string, val interface{}) error {
	var floatVal float64

	switch v := val.(type) {
	case float64:
		floatVal = v
	case int:
		floatVal = float64(v)
	case int64:
		floatVal = float64(v)
	case uint64:
		floatVal = float64(v)
	case float32:
		floatVal = float64(v)
	case bool:
		if v {
			floatVal = 1.
		} else {
			floatVal = 0.
		}
	default:
		return errors.New("error update operation: incorrect type value")
	}
	s.model.Gauge[name] = floatVal
	return nil
}

func (s *AgentService) AddCounter(name string, val int64) error {
	s.model.Counter[name] += val
	return nil
}

func (s *AgentService) Print() {
	fmt.Println(s.model)
}
