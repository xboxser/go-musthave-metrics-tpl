package service

import (
	"errors"
	"fmt"
	models "metrics/internal/model"
	"strconv"
)

// ServerService - сервис сервера отвечающий за взаимодействие с объектами метрики
type ServerService struct {
	model models.Storage
}

func NewServeService(model models.Storage) *ServerService {
	return &ServerService{
		model: model,
	}
}

// UpdateJSON - обновление метрики на основе полученного json
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

// Update - Обновление метрики по трём главным параметрам
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

// GetValueJSON - получение метрики для дальнейшей обработки в json структуре
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

// GetValue - получение метрики по типу и имени
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

// GetAll - получение информации по всем метрикам в памяти
func (s *ServerService) GetAll() map[string]string {
	gauge, counter := s.model.GetAll()

	res := make(map[string]string, len(gauge)+len(counter))
	for name, val := range gauge {
		res["gauge "+name] = strconv.FormatFloat(val, 'f', -1, 64)
	}

	for name, val := range counter {
		res["counter "+name] = strconv.FormatInt(val, 10)
	}

	return res
}

// GetModels - возвращаем все модели в нужной структуре
func (s *ServerService) GetModels() []models.Metrics {
	metrics := []models.Metrics{}
	gauge, counter := s.model.GetAll()

	for id, val := range gauge {
		metrics = append(metrics, models.Metrics{
			ID:    id,
			MType: models.Gauge,
			Value: &val,
		})
	}

	for id, val := range counter {
		metrics = append(metrics, models.Metrics{
			ID:    id,
			MType: models.Counter,
			Delta: &val,
		})
	}
	return metrics

}

// SetModel - устанавливает массив моделей в оперативную память
func (s *ServerService) SetModel(m []models.Metrics) {
	if len(m) < 1 {
		return
	}
	for _, v := range m {
		s.UpdateJSON(&v)
	}
}
