package models

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileJson struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileJson(fileName string) (*FileJson, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &FileJson{
		file:    file,
		encoder: json.NewEncoder(file),
		decoder: json.NewDecoder(file),
	}, nil
}

func (f *FileJson) Save(m []Metrics) error {
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

func (f *FileJson) Read() (*[]Metrics, error) {
	m := &[]Metrics{}
	if err := f.decoder.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func (f *FileJson) Close() error {
	return f.file.Close()
}
