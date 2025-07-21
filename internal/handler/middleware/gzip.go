package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	gzWriter *gzip.Writer
	disabled bool
}

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.disabled || w.gzWriter == nil {
		return w.ResponseWriter.Write(b)
	}
	return w.gzWriter.Write(b)
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

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		grw := &gzipWriter{ResponseWriter: w}
		defer grw.Close()

		next.ServeHTTP(grw, r)

		//
		contentType := grw.Header().Get("Content-Type")
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

		// Если тип контента не поддерживается, отменяем сжатие
		if !checkAllowedTypes {
			grw.DisableCompression()
		} else {
			grw.Header().Set("Content-Encoding", "gzip")
		}

	})
}
