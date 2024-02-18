package presenter

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/o-ga09/go122rcsample/api/internal/controller"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", controller.Health)
	mux.HandleFunc("GET /users/:id", controller.GetUsers)

	slog.Info("starting server")
	go func() {
		handler := WithTimeout(AddID(mux))
		err = http.Serve(listen, handler)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("stopping srever")
}

// AddIDはリクエスト毎にIDを付与するmiddlewareです。
func AddID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// IDを生成してcontextに保存
		id := generateID()
		ctx := context.WithValue(r.Context(), "requestId", id)
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
			fmt.Println("done")
			return
		case <-ctx.Done():
			http.Error(w, "Timeout", http.StatusRequestTimeout)
		}
	})
}

// generateIDはIDを生成するための仮の関数です。
func generateID() string {
	return uuid.NewString()
}
