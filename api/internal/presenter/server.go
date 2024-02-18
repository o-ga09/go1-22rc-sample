package presenter

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/o-ga09/go122rcsample/api/internal/controller"
	"github.com/o-ga09/go122rcsample/api/internal/middleware"
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

	reqMiddleware := middleware.AddID(mux)
	RequestLogMiddleware := middleware.RequestLogger(reqMiddleware)
	LoggerMiddleware := middleware.Logger(RequestLogMiddleware)
	handler := middleware.WithTimeout(LoggerMiddleware)

	mux.HandleFunc("GET /", controller.Health)
	mux.HandleFunc("GET /users", controller.GetUsers)

	slog.Info("starting server")
	go func() {
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
