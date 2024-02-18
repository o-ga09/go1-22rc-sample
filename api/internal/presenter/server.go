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

	HealthCheckhandler := http.HandlerFunc(controller.Health)
	HealthCheckhandler = middleware.UseMiddleware(HealthCheckhandler)

	GetUserhandler := http.HandlerFunc(controller.GetUsers)
	GetUserhandler = middleware.UseMiddleware(GetUserhandler)

	CreateUserhandler := http.HandlerFunc(controller.CreateUser)
	CreateUserhandler = middleware.UseMiddleware(CreateUserhandler)

	mux.HandleFunc("GET /", HealthCheckhandler)
	mux.HandleFunc("GET /users/{id}", GetUserhandler)
	mux.HandleFunc("POST /users", CreateUserhandler)

	slog.Info("starting server")
	go func() {
		err = http.Serve(listen, mux)
		if err != nil {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("stopping srever")
}
