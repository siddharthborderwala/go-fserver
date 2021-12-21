package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipMiddleware struct{}

func NewGzipMiddleware() *GzipMiddleware {
	return &GzipMiddleware{}
}

func (g *GzipMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			zipWriter := NewGzipResponseWriter(w)
			zipWriter.Header().Set("content-encoding", "gzip")
			next.ServeHTTP(zipWriter, r)
			defer zipWriter.Flush()
			return
		}

		next.ServeHTTP(w, r)
	})
}

type GzipResponseWriter struct {
	w      http.ResponseWriter
	zipper *gzip.Writer
}

func NewGzipResponseWriter(w http.ResponseWriter) *GzipResponseWriter {
	zipper := gzip.NewWriter(w)

	return &GzipResponseWriter{
		w,
		zipper,
	}
}

func (zipper *GzipResponseWriter) Header() http.Header {
	return zipper.w.Header()
}

func (zipper *GzipResponseWriter) Write(d []byte) (int, error) {
	return zipper.zipper.Write(d)
}

func (zipper *GzipResponseWriter) WriteHeader(statuscode int) {
	zipper.w.WriteHeader(statuscode)
}

func (zipper *GzipResponseWriter) Flush() {
	zipper.zipper.Flush()
	zipper.zipper.Close()
}
