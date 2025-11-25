package audit

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFileSubscriber(t *testing.T) {

	t.Run("checkIsNotNil", func(t *testing.T) {
		filePath := "path"
		fs := NewFileSubscriber(filePath)
		defer os.Remove(filePath)
		require.NotNil(t, fs)
		// Проверяем, что поля структуры инициализированы
		assert.NotNil(t, fs.filePath)
		assert.Equal(t, fs.filePath, filePath)
	})

	t.Run("checkIsNil", func(t *testing.T) {
		filePath := ""
		fs := NewFileSubscriber(filePath)
		defer os.Remove(filePath)
		require.NotNil(t, fs)
		// Проверяем, что поля структуры инициализированы
		assert.NotNil(t, fs.filePath)
		assert.Equal(t, fs.filePath, filePath)
	})
}
