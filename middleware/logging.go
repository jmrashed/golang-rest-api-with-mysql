package middleware

import (
	"log"
	"net/http"
	"time"
)

// LoggingResponseWriter wraps http.ResponseWriter to capture status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(data)
	lrw.size += size
	return size, err
}

// Logging middleware
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Wrap the response writer
		lrw := NewLoggingResponseWriter(w)
		
		// Process request
		next.ServeHTTP(lrw, r)
		
		// Log request details
		duration := time.Since(start)
		log.Printf(
			"%s %s %d %d bytes %v %s",
			r.Method,
			r.RequestURI,
			lrw.statusCode,
			lrw.size,
			duration,
			getClientIP(r),
		)
		
		// Log errors
		if lrw.statusCode >= 400 {
			log.Printf("ERROR: %s %s returned %d", r.Method, r.RequestURI, lrw.statusCode)
		}
	})
}