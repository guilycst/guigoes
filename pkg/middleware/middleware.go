package middleware

import (
	"compress/gzip"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
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

type GzipWriterPool struct {
	pool sync.Pool
}

func (gwp *GzipWriterPool) Get(w http.ResponseWriter) *gzip.Writer {
	gw := gwp.pool.Get()
	if gw == nil {
		var err error
		gw, err = gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			panic(err)
		}
		slog.Debug("New gzip writer created")
		return gw.(*gzip.Writer)
	}
	gw.(*gzip.Writer).Reset(w)
	return gw.(*gzip.Writer)
}

func (gwp *GzipWriterPool) Put(gw *gzip.Writer) {
	err := gw.Close()
	if err != nil {
		panic(err)
	}
	gwp.pool.Put(gw)
}

var gzPool = &GzipWriterPool{
	pool: sync.Pool{
		New: func() interface{} {
			return nil
		},
	},
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

	gzipWriter := gzPool.Get(w.ResponseWriter)
	defer gzPool.Put(gzipWriter)
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

type statusCodeExposingWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusCodeExposingWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func wrapResponseWriter(w http.ResponseWriter) http.ResponseWriter {
	return &statusCodeExposingWriter{w, http.StatusOK}
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w = wrapResponseWriter(w)
		next.ServeHTTP(w, r)
		slog.Info(fmt.Sprint(r.Method, " ", r.URL.String(), " "), "code", w.(*statusCodeExposingWriter).statusCode)
	})
}
