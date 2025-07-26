package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	gzWriter       *gzip.Writer
	isCompressible bool
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	if !w.isCompressible {
		// Если тип не поддерживается — не сжимаем
		return w.ResponseWriter.Write(b)
	}

	return w.gzWriter.Write(b)

}

func (w *gzipWriter) WriteHeader(statusCode int) {
	// Прверяем определяли ли ранее значение
	if w.isCompressible {
		return
	}

	// Проверяем Content-Type
	contentType := w.Header().Get("Content-Type")
	if contentType == "" || !isCompressible(contentType) {
		// Не сжимаем — просто отправляем оригинальный заголовок
		w.ResponseWriter.WriteHeader(statusCode)
		return
	}
	w.isCompressible = true
	// Устанавливаем заголовки для сжатия
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Del("Content-Length")

	w.ResponseWriter.WriteHeader(statusCode)
}

func isCompressible(contentType string) bool {
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

	return checkAllowedTypes
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем формат полученого запроса, если gzip то разархивируем его
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
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Запускаем сжатие ответа
		// Внутри gzipWriter есть проверка какие форматы поддерживают сжатие
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		wrapped := &gzipWriter{
			ResponseWriter: w,
			gzWriter:       gz,
		}

		// Передаём управление следующему обработчику
		next.ServeHTTP(wrapped, r)

	})
}
