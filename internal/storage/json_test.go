package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewReadJSON - тестируем readJSON
func TestNewReadJSON(t *testing.T) {
	path := "file.json"
	readJSON := NewReadJSON(path)
	require.Equal(t, path, readJSON.fileName)
}

func TestGetConfigJSON(t *testing.T) {

	tests := []struct {
		name         string
		createFile   bool
		fileContent  string
		expectError  bool
		expectedData []byte
	}{
		{
			name:         "successful read valid JSON",
			createFile:   true,
			fileContent:  `{"address": "localhost:8080", "report_interval": "10s"}`,
			expectError:  false,
			expectedData: []byte(`{"address": "localhost:8080", "report_interval": "10s"}`),
		},
		{
			name:         "successful read empty JSON object",
			createFile:   true,
			fileContent:  `{}`,
			expectError:  false,
			expectedData: []byte(`{}`),
		},
		{
			name:         "file does not exist",
			createFile:   false,
			fileContent:  "",
			expectError:  true,
			expectedData: nil,
		},
		{
			name:         "empty file",
			createFile:   true,
			fileContent:  "",
			expectError:  false,
			expectedData: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fileName string

			// создаем временный файл
			if tt.createFile {
				tmpFile, err := os.CreateTemp("", "test_config_*.json")
				require.NoError(t, err)

				fileName = tmpFile.Name()

				if tt.fileContent != "" {
					_, err = tmpFile.WriteString(tt.fileContent)
					require.NoError(t, err)
				}

				tmpFile.Close()
				defer os.Remove(fileName)
			} else {
				fileName = "nonexistent_file.json"
			}

			// создаем экземпляр ConfigJSON
			config := NewReadJSON(fileName)

			// вызываем тестируемый метод
			data, err := config.GetConfigJSON()

			// проверяем результаты
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expectedData == nil {
					assert.Nil(t, data)
				} else {
					assert.Equal(t, tt.expectedData, data)
				}
			}
		})
	}
}
