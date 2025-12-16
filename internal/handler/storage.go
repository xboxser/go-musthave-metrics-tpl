package handler

import (
	"context"
	"metrics/internal/config"
	"metrics/internal/config/db"
	models "metrics/internal/model"
	"metrics/internal/storage"

	"go.uber.org/zap"
)

// Storage - хранение и чтение метрик
type StorageManager struct {
	config *config.ConfigServer
	file   *storage.FileJSON
	db     *db.DB
	sugar  zap.SugaredLogger
}

func NewStorageManager(config *config.ConfigServer) *StorageManager {
	var sugar zap.SugaredLogger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar = *logger.Sugar()

	file, err := storage.NewFileJSON(config.FileStoragePath)
	if err != nil {
		panic(err)
	}
	return &StorageManager{
		config: config,
		file:   file,
		sugar:  sugar,
	}
}

func (h *StorageManager) Read() []models.Metrics {
	if !h.config.Restore {
		return nil
	}
	models, ok := h.ReadFromDB()
	if ok {
		return models
	}

	return h.ReadFromFile()
}

func (h *StorageManager) ReadFromFile() []models.Metrics {
	m, err := h.file.Read()
	if err != nil {
		panic(err)
	}

	return *m
}

func (h *StorageManager) SaveToFile(m []models.Metrics) error {
	err := h.file.Save(m)
	return err
}

func (h *StorageManager) Close() error {
	return h.file.Close()
}

func (h *StorageManager) ConnectDB(ctx context.Context) error {
	if h.config.DateBaseDSN == "" {
		return nil
	}
	db, err := db.NewDB(ctx, h.config.DateBaseDSN)
	if err != nil {
		return err
	}
	h.db = db
	return nil
}

func (h *StorageManager) SaveToDB(models []models.Metrics) bool {
	if h.config.DateBaseDSN == "" || !h.db.Ping() {
		return false
	}
	err := h.db.SaveAll(models)
	if err != nil {
		h.sugar.Errorf("Ошибка при записи в БД: %v\n", err)
	} else {
		return true
	}

	return false
}

func (h *StorageManager) ReadFromDB() ([]models.Metrics, bool) {
	if h.config.DateBaseDSN == "" || !h.db.Ping() {
		return nil, false
	}

	m, err := h.db.ReadAll()
	if err != nil {
		h.sugar.Errorln("Не удалось получить информацию из БД", err)
		return nil, false
	}

	return m, true
}

// Ping - проверка есть ли подключение к БД
func (h *StorageManager) Ping() bool {
	if h.config.DateBaseDSN == "" {
		return false
	}
	return h.db.Ping()
}
