package storage

import (
	"bufio"
	"os"
)

// ReadJSON - читает данные из json файла
type ReadJSON struct {
	fileName string
}

func NewReadJSON(fileName string) *ReadJSON {
	return &ReadJSON{fileName: fileName}
}

func (c *ReadJSON) GetConfigJSON() ([]byte, error) {
	file, err := os.Open(c.fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return nil, nil
	}

	bufioReader := bufio.NewReader(file)
	fileContent := make([]byte, fileSize)
	_, err = bufioReader.Read(fileContent)
	if err != nil {
		return nil, err
	}

	return fileContent, nil
}
