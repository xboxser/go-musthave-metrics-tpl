package middleware

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	disabled bool
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if w.disabled || w.gzWriter == nil {
		return w.ResponseWriter.Write(b)
	}

	// Только при первом Write() проверяем Content-Type
	if !w.disabled && w.gzWriter != nil {
		contentType := w.Header().Get("Content-Type")
		if !isCompressible(contentType) {
			w.DisableCompression()
			return w.ResponseWriter.Write(b)
		}

		// Устанавливаем заголовки сжатия
		if w.Header().Get("Content-Encoding") == "" {
			w.Header().Del("Content-Length")
			w.Header().Set("Content-Encoding", "gzip")
		}
	}

	return w.gzWriter.Write(b)
}

func isCompressible(contentType string) bool {
	fmt.Println("contentType", contentType)
	allowedTypes := []string{
		"application/json",
		"text/html",
	}

	checkAllowedTypes := false
	for _, v := range allowedTypes {
		if strings.Contains(contentType, v) {
			checkAllowedTypes = true
			break
		}
	}
	fmt.Println("checkAllowedTypes", checkAllowedTypes)

	return checkAllowedTypes
}

func (w *gzipWriter) Close() {
	if w.gzWriter != nil {
		w.gzWriter.Close()
	}
}

func (w *gzipWriter) DisableCompression() {
	w.disabled = true
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Failed to decompress request", http.StatusBadRequest)
				return
			}
			defer gz.Close()
			r.Body = gz
		}

		// проверяем ждет ли ответа в формате gzip
		fmt.Println("Accept-Encoding", r.Header.Get("Accept-Encoding"))
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Создаём обёртку
		gzWriter := gzip.NewWriter(w)
		grw := &gzipWriter{
			ResponseWriter: w,
			gzWriter:       gzWriter,
			disabled:       false,
		}

		defer grw.Close()

		// Передаём управление следующему обработчику
		next.ServeHTTP(grw, r)

	})
}
