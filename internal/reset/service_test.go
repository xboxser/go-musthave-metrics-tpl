package reset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceReset(t *testing.T) {
	root := "root"
	service := NewServiceReset(root)
	// Проверяем, что структура создана
	assert.NotNil(t, service)

	// Проверяем, что поля структуры инициализированы
	assert.NotNil(t, service.root)
	assert.NotNil(t, service.PackageInfo)
	assert.Equal(t, service.root, root)
}
