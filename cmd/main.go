package main

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hellow World go 1.22 ! from GET\n")
	})
	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hellow World go 1.22 ! from POST\n")
	})

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
