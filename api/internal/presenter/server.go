package presenter

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
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
	mux.HandleFunc("GET /", health)
	mux.HandleFunc("GET /users/:id", getUsers)

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

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hellow World go 1.22 ! from GET\n")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sql := "SELECT * FROM users WHERE id = ?"
	rows, err := db.Query(sql, id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var email string
		err = rows.Scan(&id, &name, &email)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "id: %d, name: %s, email: %s\n", id, name, email)
	}
}
