package rest_errors

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

var (
	ErrInternalServer   = errors.New("internal server error")
	ErrBadRequest       = errors.New("bad request")
	ErrUserUnauthorized = errors.New("user unauthorized")
)

type ResponseError struct {
	Error string `json:"error"`
}

func (re *ResponseError) String() string {
	jsonData, _ := json.Marshal(re)
	return string(jsonData)
}

func HandleError(w http.ResponseWriter, err error, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := &ResponseError{Error: err.Error()}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("Error encoding error response", "error", err)
		http.Error(w, `{"error": "internal error"}`, http.StatusInternalServerError)
	}
}
