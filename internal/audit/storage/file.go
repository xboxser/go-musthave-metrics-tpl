package storage

import (
	"encoding/json"
	"metrics/internal/audit/model"
	"os"
)

// FileAuditJSON - файл для хранения аудита в json
type FileAuditJSON struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewAuditFileJSON(fileName string) (*FileAuditJSON, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Настроим один раз

	return &FileAuditJSON{
		file:    file,
		encoder: encoder,
		decoder: json.NewDecoder(file),
	}, nil
}

// Save - сохранение аудита в конец файла
func (f *FileAuditJSON) Save(a model.Audit) error {
	err := f.encoder.Encode(a)
	if err != nil {
		return err
	}

	return nil
}

func (f *FileAuditJSON) Close() error {
	return f.file.Close()
}
