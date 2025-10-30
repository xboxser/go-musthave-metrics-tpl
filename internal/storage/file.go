package storage

import (
	"encoding/json"
	models "metrics/internal/model"
	"os"
)

// Сохраняет метрики в файл
type FileJSON struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileJSON(fileName string) (*FileJSON, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Настроим один раз

	return &FileJSON{
		file:    file,
		encoder: encoder,
		decoder: json.NewDecoder(file),
	}, nil
}

func (f *FileJSON) Save(m []models.Metrics) error {
	if err := f.file.Truncate(0); err != nil {
		return err
	}

	if _, err := f.file.Seek(0, 0); err != nil {
		return err
	}

	err := f.encoder.Encode(m)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileJSON) Read() (*[]models.Metrics, error) {
	m := &[]models.Metrics{}
	// Проверяем размер файла
	stat, err := f.file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() == 0 {
		return m, nil
	}

	// Пересоздаем decoder
	f.decoder = json.NewDecoder(f.file)

	if err := f.decoder.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *FileJSON) Close() error {
	return f.file.Close()
}
