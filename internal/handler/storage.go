package handler

import (
	"metrics/internal/config"
	models "metrics/internal/model"
	"metrics/internal/storage"
)

type Storage struct {
	config *config.ConfigServer
	file   *storage.FileJSON
}

func NewStorage(config *config.ConfigServer) *Storage {
	file, err := storage.NewFileJSON(config.FileStoragePath)
	if err != nil {
		panic(err)
	}
	return &Storage{
		config: config,
		file:   file,
	}
}

func (h *Storage) ReadFromFile() []models.Metrics {
	m, err := h.file.Read()
	if err != nil {
		panic(err)
	}

	return *m
}

func (h *Storage) SaveToFile(m []models.Metrics) error {
	err := h.file.Save(m)
	return err
}

func (h *Storage) Close() error {
	return h.file.Close()
}
