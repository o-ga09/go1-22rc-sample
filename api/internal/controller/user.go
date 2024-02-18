package controller

import (
	"database/sql"
	"encoding/json"
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
		slog.Log(r.Context(), middleware.SeverityError, "con not get environment value")
		return
	}
	defer func() {
		slog.Log(r.Context(), middleware.SeverityInfo, "db disconnect ....", "requestId", middleware.GetRequestID(r.Context()))
		db.Close()
	}()

	query := "SELECT * FROM users WHERE id = ?"
	var uid int
	var name string
	var email string
	err = db.QueryRow(query, id).Scan(&uid, &email, &name)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Log(r.Context(), middleware.SeverityInfo, "no rows", "requestId", middleware.GetRequestID(r.Context()))
			return
		}
		slog.Log(r.Context(), middleware.SeverityError, "panic error", "error message", err, "requestId", middleware.GetRequestID(r.Context()))
	}
	result := struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    uid,
		Name:  name,
		Email: email,
	}
	slog.Log(r.Context(), middleware.SeverityInfo, "result", "data", result, "requestId", middleware.GetRequestID(r.Context()))
	middleware.Response(&w, http.StatusOK, result)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Log(r.Context(), middleware.SeverityInfo, "can not get request body", "requestId", middleware.GetRequestID(r.Context()))
		return
	}

	cfg, _ := config.New()
	db, err := sql.Open("mysql", cfg.Database_url)
	if err != nil {
		slog.Log(r.Context(), middleware.SeverityError, "db connect error...")
		panic(err)
	}
	defer func() {
		slog.Log(r.Context(), middleware.SeverityInfo, "db disconnect ....", "requestId", middleware.GetRequestID(r.Context()))
		db.Close()
	}()

	name := reqBody.Name
	email := reqBody.Email
	sql := "INSERT INTO users (name, email) VALUES (?, ?)"
	_, err = db.Exec(sql, name, email)
	if err != nil {
		slog.Log(r.Context(), middleware.SeverityError, "can not insert", "error message", err, "requestId", middleware.GetRequestID(r.Context()))
		return
	}
	middleware.Response(&w, http.StatusCreated, reqBody)
}
