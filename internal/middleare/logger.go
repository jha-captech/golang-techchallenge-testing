package middleare

import (
	"log/slog"
	"net/http"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

// Logger is a middleware that logs the request method, path, duration, and
// status code.
func Logger(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(wrapped, r)

			logger.InfoContext(
				r.Context(),
				"request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("duration", time.Since(start).String()),
				slog.Int("status", wrapped.statusCode),
			)
		})
	}
}
