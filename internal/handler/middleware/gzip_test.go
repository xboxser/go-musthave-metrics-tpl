package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestIsCompressible - проверяет корректность работы isCompressible
func TestIsCompressible(t *testing.T) {
	type args struct {
		contentType string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "valid parameters 1", args: args{contentType: "application/json"}, want: true},
		{name: "valid parameters 2", args: args{contentType: "text/html"}, want: true},
		{name: "no valid parameters", args: args{contentType: "zip"}, want: false},
		{name: "empty parameters", args: args{contentType: ""}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isCompressible(tt.args.contentType)
			if result != tt.want {
				t.Errorf("isCompressible(%q) = %v; expected %v", tt.args.contentType, result, tt.want)
			}
		})
	}
}

// TestGzipMiddleware_NoCompressionWithoutHeader - проверяет отсутствие сжатия без заголовка
func TestGzipMiddleware_NoCompressionWithoutHeader(t *testing.T) {
	testData := "simple test data"
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(testData))
	})

	gzipHandler := GzipMiddleware(handler)

	// Запрос БЕЗ Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()

	gzipHandler.ServeHTTP(rec, req)

	// Не должно быть заголовка Content-Encoding
	if rec.Header().Get("Content-Encoding") != "" {
		t.Error("Content-Encoding should be empty without Accept-Encoding")
	}

	// Данные должны быть несжатыми
	if rec.Body.String() != testData {
		t.Error("response body should be uncompressed")
	}
}

// TestGzipMiddleware_CompressResponse - проверяет сжатие ответа
func TestGzipMiddleware_CompressResponse(t *testing.T) {
	// создаём тестовый handler, который возвращает большой текст
	testData := strings.Repeat("test data for compression ", 100)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testData))
	})

	// Оборачиваем в gzip middleware
	gzipHandler := GzipMiddleware(handler)

	// Создаём запрос с Accept-Encoding: gzip
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	gzipHandler.ServeHTTP(rec, req)

	// Проверяем статус
	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// Проверяем заголовок Content-Encoding
	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Error("Content-Encoding header should be 'gzip'")
	}

	// Проверяем, что данные действительно сжаты
	body := rec.Body.Bytes()

	// Распаковываем для проверки
	gzReader, err := gzip.NewReader(bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	decompressed, err := io.ReadAll(gzReader)
	if err != nil {
		t.Fatalf("failed to decompress: %v", err)
	}

	if string(decompressed) != testData {
		t.Error("decompressed data doesn't match original")
	}

	// Проверяем, что сжатие действительно произошло
	if len(body) >= len(testData) {
		t.Error("compressed data should be smaller than original")
	}
}

// TestGzipMiddleware_Headers - проверяет корректность заголовков
func TestGzipMiddleware_Headers(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Custom-Header", "test-value")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"test":"data"}`))
	})

	gzipHandler := GzipMiddleware(handler)

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	rec := httptest.NewRecorder()
	gzipHandler.ServeHTTP(rec, req)

	// Проверяем, что кастомные заголовки сохранились
	if rec.Header().Get("X-Custom-Header") != "test-value" {
		t.Error("custom headers should be preserved")
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Error("content-type should be preserved")
	}
}
