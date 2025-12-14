package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CIDRValidate - проверяет входит ли указанный адрес в нужную подсеть
func TestCIDRValidate(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name      string
		cidr      string
		ipAddress string
		want      bool
		wantErr   bool
	}{
		{name: "valid parameters 1", cidr: "192.168.1.0/24", ipAddress: "192.168.1.10", want: true, wantErr: false},
		{name: "valid parameters 2", cidr: "192.168.1.0/24", ipAddress: "192.168.1.25", want: true, wantErr: false},
		{name: "IP address outside subnet", cidr: "192.168.1.0/24", ipAddress: "192.168.2.10", want: false, wantErr: false},
		{name: "IP address at network boundary", cidr: "10.0.0.0/8", ipAddress: "10.0.0.1", want: true, wantErr: false},
		{name: "IP address at broadcast boundary", cidr: "10.0.0.0/8", ipAddress: "10.255.255.254", want: true, wantErr: false},
		{name: "Invalid CIDR format", cidr: "192.168.1.0/33", ipAddress: "192.168.1.10", want: false, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := CIDRValidate(tt.cidr, tt.ipAddress)

			require.Equal(t, ok, tt.want)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCIDRMiddleware(t *testing.T) {
	// Создаем тестовый обработчик, который будет вызван, если middleware пропустит запрос
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name           string
		trustedSubnet  string
		xRealIP        string
		expectedStatus int
		description    string
	}{
		{
			name:           "Empty trusted subnet",
			trustedSubnet:  "",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusOK,
			description:    "Когда подсеть не указана, запрос должен пройти",
		},
		{
			name:           "Valid IP in trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusOK,
			description:    "IP из доверенной подсети должен пройти",
		},
		{
			name:           "Invalid IP outside trusted subnet",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "10.0.0.1",
			expectedStatus: http.StatusForbidden,
			description:    "IP вне доверенной подсети должен быть заблокирован",
		},
		{
			name:           "Invalid subnet format",
			trustedSubnet:  "192.168.1.0/33",
			xRealIP:        "192.168.1.100",
			expectedStatus: http.StatusForbidden,
			description:    "Неверный формат подсети должен блокировать запрос",
		},
		{
			name:           "No X-Real-IP header",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "",
			expectedStatus: http.StatusForbidden,
			description:    "Отсутствие заголовка X-Real-IP должен блокировать запрос",
		},
		{
			name:           "Invalid IP format",
			trustedSubnet:  "192.168.1.0/24",
			xRealIP:        "invalid.ip.address",
			expectedStatus: http.StatusForbidden,
			description:    "Неверный формат IP должен блокировать запрос",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем middleware с параметрами теста
			middleware := CIDRMiddleware(tt.trustedSubnet)

			// Оборачиваем наш тестовый обработчик в middleware
			handler := middleware(nextHandler)

			// Создаем тестовый HTTP запрос
			req := httptest.NewRequest("GET", "/", nil)
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}

			// Создаем рекордер для записи ответа
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			handler.ServeHTTP(rr, req)

			// Проверяем статус ответа
			assert.Equal(t, tt.expectedStatus, rr.Code, tt.description)
		})
	}
}
