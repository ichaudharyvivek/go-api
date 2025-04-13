package httpx

import (
	"encoding/json"
	"net/http"
)

// APIResponse is the general response structure
type APIResponse struct {
	Success bool     `json:"success"`
	Data    any      `json:"data,omitempty"`
	Error   string   `json:"error,omitempty"`
	Errors  []string `json:"errors,omitempty"`
}

// WriteJSON writes a fully customized APIResponse
func respondJSON(w http.ResponseWriter, status int, payload APIResponse) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

// Success returns 200 OK with data
func Ok(w http.ResponseWriter, data any) error {
	return respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Created writes a success response with data and `201` status
func Created(w http.ResponseWriter, data any) error {
	return respondJSON(w, http.StatusCreated, APIResponse{
		Success: true,
		Data:    data,
	})
}

// Error writes a single error message
func Error(w http.ResponseWriter, msg string, status int) {
	respondJSON(w, status, APIResponse{
		Success: false,
		Error:   msg,
	})
}

// Writes a list of multiple errors
func Errors(w http.ResponseWriter, errs []string, status int) {
	respondJSON(w, status, APIResponse{
		Success: false,
		Errors:  errs,
	})
}
