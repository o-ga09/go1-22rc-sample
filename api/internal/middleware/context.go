package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/o-ga09/go122rcsample/api/pkg"
)

type RequestId string

// AddIDはリクエスト毎にIDを付与するmiddlewareです。
func AddID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IDを生成してcontextに保存
		id := pkg.GenerateID()
		ctx := context.WithValue(r.Context(), RequestId("requestId"), id)
		// 次のハンドラーに渡す
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithTimeoutはIDを追加するmiddlewareです。
func WithTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Context() == nil {
			r = r.WithContext(context.Background())
		}

		// タイムアウトを設定
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel() // 処理が終了したらキャンセルする

		// 次のハンドラーを実行し、タイムアウトが発生した場合はエラーメッセージを出力
		done := make(chan struct{})
		go func() {
			defer close(done)
			next.ServeHTTP(w, r)
		}()
		select {
		case <-done:
			// ハンドラーが正常に終了した場合は何もしない
			return
		case <-ctx.Done():
			http.Error(w, "Timeout", http.StatusRequestTimeout)
		}
	})
}

func GetRequestID(ctx context.Context) string {
	return ctx.Value(RequestId("requestId")).(string)
}
