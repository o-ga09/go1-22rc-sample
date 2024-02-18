package middleware

import (
	"encoding/json"
	"net/http"
)

func Response(w *http.ResponseWriter, status int, message interface{}) {
	json, err := json.Marshal(message)
	if err != nil {
		(*w).Header().Set("Content-Type", "application/json")
		(*w).WriteHeader(http.StatusInternalServerError)
		(*w).Write([]byte(`{"message": "error marshalling json"}`))
		return
	}

	(*w).Header().Set("Content-Type", "application/json")
	(*w).WriteHeader(status)
	(*w).Write(json)
}
