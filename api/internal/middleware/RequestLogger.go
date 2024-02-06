package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type RequestInfo struct {
	status         int
	contentsLength int64
	path           string
	sourceIP       string
	query          string
	userAgent      string
	errors         string
	elapsed        time.Duration
}

func (r *RequestInfo) LogValue() interface{} { // Assuming slog expects an interface{}
	return map[string]interface{}{
		"status":          r.status,
		"contents_length": r.contentsLength,
		"path":            r.path,
		"sourceIP":        r.sourceIP,
		"query":           r.query,
		"user_agent":      r.userAgent,
		"errors":          r.errors,
		"elapsed":         r.elapsed.String(),
	}
}

func RequestLogger(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			defer func() {
				r := &RequestInfo{
					status:         r.Response.StatusCode, // Get status after request handling
					contentsLength: r.ContentLength,
					path:           r.URL.Path,
					sourceIP:       r.RemoteAddr,
					query:          r.URL.RawQuery,
					userAgent:      r.UserAgent(),
					errors:// Implement error handling logic if needed
					"",
					elapsed: time.Since(start),
				}
				slog.Log(context.TODO(), SeverityError, "Request Info", "Request", r.LogValue()) // Adjust logging context as needed
			}()

			next.ServeHTTP(w, r)
		})
	}
}
