package service

import (
	"errors"
	"fmt"
	models "metrics/internal/model"
	"strconv"
)

type ServerService struct {
	model models.Storage
}

func NewServeService(model models.Storage) *ServerService {
	return &ServerService{
		model: model,
	}
}

func (s *ServerService) UpdateJSON(m *models.Metrics) error {
	if m.MType == models.Counter {
		s.model.UpdateCounter(m.ID, *m.Delta)
		val, _ := s.model.GetCounter(m.ID)
		m.Delta = &val
	} else if m.MType == models.Gauge {
		s.model.UpdateGauge(m.ID, *m.Value)
		val, _ := s.model.GetGauge(m.ID)
		m.Value = &val
	} else {
		return errors.New("error update operation: incorrect type")
	}
	return nil
}

func (s *ServerService) Update(t string, name string, val string) error {

	if t == models.Counter {
		val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.model.UpdateCounter(name, val)
	} else if t == models.Gauge {
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.model.UpdateGauge(name, val)
	} else {
		return errors.New("error update operation: incorrect type")
	}
	return nil
}

func (s *ServerService) GetValueJSON(m *models.Metrics) error {
	if m.MType == models.Counter {
		val, ok := s.model.GetCounter(m.ID)
		if !ok {
			return errors.New("empty value")
		}
		m.Delta = &val
	} else if m.MType == models.Gauge {
		val, ok := s.model.GetGauge(m.ID)
		if !ok {
			return errors.New("empty value")
		}
		m.Value = &val
	} else {
		return errors.New("error update operation: incorrect type")
	}
	return nil
}

func (s *ServerService) GetValue(t string, name string) (string, error) {
	if t == models.Counter {
		val, ok := s.model.GetCounter(name)
		if !ok {
			return "", errors.New("empty value")
		}
		return strconv.FormatInt(val, 10), nil
	} else if t == models.Gauge {
		val, ok := s.model.GetGauge(name)
		if !ok {
			return "", errors.New("empty value")
		}
		return fmt.Sprintf("%g", val), nil

	} else {
		return "", errors.New("error get operation: incorrect type")
	}
}

func (s *ServerService) GetAll() map[string]string {
	guuge, counter := s.model.GetAll()

	res := map[string]string{}
	for name, val := range guuge {
		res["guuge "+name] = fmt.Sprintf("%g", val)
	}

	for name, val := range counter {
		res["counter "+name] = fmt.Sprintf("%v", val)
	}

	return res
}
