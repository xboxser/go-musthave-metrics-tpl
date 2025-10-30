package storage

import (
	"encoding/json"
	"metrics/internal/audit/model"
	"os"
)

type FileAuditJSON struct {
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewAuditFileJSON(fileName string) (*FileAuditJSON, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0666)
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

func (f *FileAuditJSON) Save(a model.Audit) error {
	err := f.encoder.Encode(a)
	if err != nil {
		return err
	}

	return nil
}
