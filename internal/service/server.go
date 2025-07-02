package service

import (
	"errors"
	"fmt"
	models "metrics/internal/model"
	"strconv"
)

type ServerService struct {
	model *models.MemStorage
}

func NewServeService(model *models.MemStorage) *ServerService {
	return &ServerService{
		model: model,
	}
}

func (s *ServerService) Update(t string, name string, val string) error {

	if t == models.Counter {
		val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.model.Counter[name] += val
	} else if t == models.Gauge {
		val, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return errors.New("error update operation: incorrect type value")
		}
		s.model.Gauge[name] = val
	} else {
		return errors.New("error update operation: incorrect type")
	}
	fmt.Println(s)
	return nil

}
