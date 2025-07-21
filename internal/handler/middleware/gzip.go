package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(b))
	}
	contentType := w.Header().Get("Content-Type")
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	allowedTypes := []string{
		"application/json",
		"text/html",
	}

	checkAllowedTypes := false
	for _, v := range allowedTypes {

		if strings.Contains(contentType, v) {
			checkAllowedTypes = true
		}
	}
	// тип контента не поддерживается, не сжимаем
	if !checkAllowedTypes {
		return w.ResponseWriter.Write(b)
	}

	return w.Writer.Write(b)
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

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gzw := gzip.NewWriter(w)
		defer gzw.Close()

		w.Header().Set("Content-Encoding", "gzip")

		// Обертка для ResponseWriter для перехвата записи
		grw := &gzipWriter{
			ResponseWriter: w,
			Writer:         gzw,
		}

		next.ServeHTTP(grw, r)

	})
}
