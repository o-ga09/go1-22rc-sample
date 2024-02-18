package controller

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("requestId").(string)
	fmt.Printf("Health: %s\n", id)
	fmt.Fprint(w, "Hellow World go 1.22 ! from GET\n")
}
