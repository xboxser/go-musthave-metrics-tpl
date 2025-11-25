package audit

import (
	"metrics/internal/audit/model"
	"metrics/internal/audit/storage"
)

// FileSubscriber - обработчик аудита работающий  с файлом  JSON
// Реализует интерфейс Observer
// файла-приёмник, добавляет информацию в конец файла, указанного в параметре конфигурации, на новой строке.
type FileSubscriber struct {
	file     *storage.FileAuditJSON
	filePath string
}

func NewFileSubscriber(filePath string) *FileSubscriber {

	if filePath == "" {
		return &FileSubscriber{filePath: filePath}
	}
	file, err := storage.NewAuditFileJSON(filePath)
	if err != nil {
		panic(err)
	}
	return &FileSubscriber{
		filePath: filePath,
		file:     file,
	}
}
func (f *FileSubscriber) Update(audit model.Audit) {
	if f.file != nil {
		err := f.file.Save(audit)
		if err != nil {
			f.file.Close()
			panic(err)
		}
	}
}
