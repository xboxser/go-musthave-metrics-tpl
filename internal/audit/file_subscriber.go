package audit

import (
	"fmt"
	"metrics/internal/audit/model"
	"metrics/internal/audit/storage"
)

// Реализует интерфейс Observer
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
	fmt.Printf("The %s FileSubscriber is notified of the %v event\n", f.filePath, audit)
	if f.file != nil {
		err := f.file.Save(audit)
		if err != nil {
			panic(err)
		}
	}
}
