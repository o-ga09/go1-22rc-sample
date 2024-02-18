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
	HealthCheckhandler = middleware.WithTimeout(HealthCheckhandler)
	HealthCheckhandler = middleware.RequestLogger(HealthCheckhandler)
	HealthCheckhandler = middleware.AddID(HealthCheckhandler)
	HealthCheckhandler = middleware.Logger(HealthCheckhandler)

	Userhandler := http.HandlerFunc(controller.GetUsers)
	Userhandler = middleware.WithTimeout(Userhandler)
	Userhandler = middleware.RequestLogger(Userhandler)
	Userhandler = middleware.AddID(Userhandler)
	Userhandler = middleware.Logger(Userhandler)

	mux.HandleFunc("/", HealthCheckhandler)
	mux.Handle("/users", Userhandler)

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
