package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type RequestInfo struct {
	ContentsLength int64
	Path           string
	SourceIP       string
	Query          string
	UserAgent      string
	Errors         string
	Elapsed        time.Duration
}

func (r *RequestInfo) LogValue() interface{} { // Assuming slog expects an interface{}
	return map[string]interface{}{
		"contents_length": r.ContentsLength,
		"path":            r.Path,
		"sourceIP":        r.SourceIP,
		"query":           r.Query,
		"user_agent":      r.UserAgent,
		"errors":          r.Errors,
		"elapsed":         r.Elapsed.String(),
	}
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Log(r.Context(), SeverityInfo, "処理開始", "requestId", r.Context())
		start := time.Now()

		next.ServeHTTP(w, r)

		req := RequestInfo{
			ContentsLength: r.ContentLength,
			Path:           r.RequestURI,
			SourceIP:       r.RemoteAddr,
			Query:          r.URL.RawQuery,
			UserAgent:      r.UserAgent(),
			Errors:         "errors",
			Elapsed:        time.Since(start),
		}

		slog.Log(r.Context(), SeverityInfo, "処理終了", "Request", req.LogValue()) // Adjust logging context as needed
	})
}
