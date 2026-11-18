package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "erro ao gerar resposta", http.StatusInternalServerError)
	}
}

func ErrorResponseJSON(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, ErrorResponse{
		Error: message,
	})
}