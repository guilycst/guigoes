package middleware

import (
	"compress/gzip"
	"log/slog"
	"net/http"
	"strings"
)

func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "private, max-age=31536000")
		next.ServeHTTP(w, r)
	})
}

func PanicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//recover from panic
		defer func() {
			if r := recover(); r != nil {
				slog.Error("Recovered from panic:", r)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	*http.Request
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	ct := w.Header().Get("Content-Type")
	if ct == "" {
		ct = http.DetectContentType(b)
	}

	if !strings.Contains(ct, "text") && !strings.Contains(ct, "json") && !strings.Contains(ct, "javascript") {
		return w.ResponseWriter.Write(b)
	}

	w.Header().Set("Content-Encoding", "gzip")

	gzipWriter, _ := gzip.NewWriterLevel(w.ResponseWriter, gzip.BestSpeed)
	defer gzipWriter.Close()

	return gzipWriter.Write(b)
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		w = &gzipResponseWriter{ResponseWriter: w, Request: r}
		next.ServeHTTP(w, r)
	})
}
