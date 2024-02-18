package controller

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/o-ga09/go122rcsample/api/internal/config"
	"github.com/o-ga09/go122rcsample/api/internal/middleware"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	cfg, _ := config.New()
	db, err := sql.Open("mysql", cfg.Database_url)
	if err != nil {
		slog.Log(r.Context(), middleware.SeverityError, "db connect error...")
		panic(err)
	}
	defer db.Close()
	slog.Log(r.Context(), middleware.SeverityInfo, "db connect ...", "requestId", middleware.GetRequestID(r.Context()))

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
