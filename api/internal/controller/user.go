package controller

import (
	"database/sql"
	"fmt"
	"net/http"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
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
