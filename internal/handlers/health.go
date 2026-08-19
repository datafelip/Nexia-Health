package handlers

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}