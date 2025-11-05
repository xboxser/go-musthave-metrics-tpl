package sender

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewSender(t *testing.T) {
	baseURL := "localhost:8080"
	sender := NewSender(&baseURL)

	if sender == nil {
		t.Fatal("NewSender returned nil")
	}
	if sender.baseURL != &baseURL {
		t.Errorf("expected baseURL %s", baseURL)
	}
	if sender.client == nil {
		t.Error("client should not be nil")
	}
	if &sender.sugar == nil {
		t.Error("sugar should not be nil")
	}
}

// TestSender_Send_WithHash проверяет отправку с хешем
func TestSender_Send_WithHash(t *testing.T) {
	expectedHash := "test_hash_123"

	//Создаем тестовый сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем наличие хеша в заголовке
		hash := r.Header.Get("HashSHA256")
		if hash != expectedHash {
			t.Errorf("expected hash %s, got %s", expectedHash, hash)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Получить адрес сервера без протокола
	baseURL := server.URL[7:]
	sender := NewSender(&baseURL)

	var compressedBuf bytes.Buffer
	statusCode, err := sender.Send(compressedBuf, expectedHash)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if statusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, statusCode)
	}
}

// TestSender_SendRequest_Success проверяет успешную отправку через SendRequest
func TestSender_SendRequest_Success(t *testing.T) {
	requestCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	baseURL := server.URL[7:]
	sender := NewSender(&baseURL)

	json := []byte(`{"metric":"counter","value":100}`)
	err := sender.SendRequest(json)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if requestCount != 1 {
		t.Errorf("expected 1 request, got %d", requestCount)
	}
}

// TestSender_Send_InvalidURL проверяет обработку невалидного URL
func TestSender_Send_InvalidURL(t *testing.T) {
	// Создаём невалидный baseURL
	invalidURL := "://invalid url with spaces"
	sender := NewSender(&invalidURL)

	var compressedBuf bytes.Buffer
	statusCode, err := sender.Send(compressedBuf, "")

	if err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
	if statusCode != 0 {
		t.Errorf("expected status 0 for error, got %d", statusCode)
	}
}

// TestSender_Send_ConnectionRefused проверяет обработку недоступного сервера
func TestSender_Send_ConnectionRefused(t *testing.T) {
	// Используем порт, на котором точно никто не слушает
	baseURL := "localhost:19999"
	sender := NewSender(&baseURL)

	var compressedBuf bytes.Buffer
	statusCode, err := sender.Send(compressedBuf, "")

	if err == nil {
		t.Error("expected error for connection refused, got nil")
	}
	if statusCode != 0 {
		t.Errorf("expected status 0 for connection error, got %d", statusCode)
	}
}

// TestSender_SendRequest_AllRetriesFailed проверяет исчерпание всех попыток
func TestSender_SendRequest_AllRetriesFailed(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping retry test in short mode")
	}

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	baseURL := server.URL[7:]
	sender := NewSender(&baseURL)

	json := []byte(`{"test":"fail"}`)
	err := sender.SendRequest(json)

	if err == nil {
		t.Error("expected error when all retries fail, got nil")
	}
	// Должно быть 4 попытки (0s, 1s, 3s, 5s)
	if requestCount != 4 {
		t.Errorf("expected 4 retry attempts, got %d", requestCount)
	}
}
