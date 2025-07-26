package storage

import (
	"encoding/json"
	"fmt"
	models "metrics/internal/model"
	"os"
)

type FileJSON struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileJSON(fileName string) (*FileJSON, error) {
	fmt.Println("fileName", fileName)
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileJSON{
		file:    file,
		encoder: json.NewEncoder(file),
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
	f.encoder.SetIndent("", "  ") // Форматирование с отступами
	err := f.encoder.Encode(m)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return err
	}

	fmt.Println("Данные успешно сохранены")
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

	if err := f.decoder.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *FileJSON) Close() error {
	return f.file.Close()
}
