package controller

import (
	"fmt"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hellow World go 1.22 ! from GET\n")
}
